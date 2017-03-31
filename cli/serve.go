package cli

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/e-dard/netbug"
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/kolide/kolide/server/config"
	"github.com/kolide/kolide/server/datastore/mysql"
	"github.com/kolide/kolide/server/kolide"
	"github.com/kolide/kolide/server/license"
	"github.com/kolide/kolide/server/mail"
	"github.com/kolide/kolide/server/pubsub"
	"github.com/kolide/kolide/server/service"
	"github.com/kolide/kolide/server/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type initializer interface {
	// Initialize is used to populate a datastore with
	// preloaded data
	Initialize() error
}

func createServeCmd(configManager config.Manager) *cobra.Command {
	// Whether to enable the debug endpoints
	debug := false

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Launch the kolide server",
		Long: `
Launch the kolide server

Use kolide serve to run the main HTTPS server. The Kolide server bundles
together all static assets and dependent libraries into a statically linked go
binary (which you're executing right now). Use the options below to customize
the way that the kolide server works.
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := configManager.LoadConfig()

			var logger kitlog.Logger
			{
				output := os.Stderr
				if config.Logging.JSON {
					logger = kitlog.NewJSONLogger(output)
				} else {
					logger = kitlog.NewLogfmtLogger(output)
				}
				logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
			}

			var ds kolide.Datastore
			var err error
			mailService := mail.NewService()

			ds, err = mysql.New(config.Mysql, clock.C, mysql.Logger(logger))
			if err != nil {
				initFatal(err, "initializing datastore")
			}

			migrationStatus, err := ds.MigrationStatus()
			if err != nil {
				initFatal(err, "retrieving migration status")
			}

			switch migrationStatus {
			case kolide.SomeMigrationsCompleted:
				fmt.Printf("################################################################################\n"+
					"# WARNING:\n"+
					"#   Your Kolide database is missing required migrations. This is likely to cause\n"+
					"#   errors in Kolide.\n"+
					"#\n"+
					"#   Run `%s prepare db` to perform migrations.\n"+
					"################################################################################\n",
					os.Args[0])

			case kolide.NoMigrationsCompleted:
				fmt.Printf("################################################################################\n"+
					"# ERROR:\n"+
					"#   Your Kolide database is not initialized. Kolide cannot start up.\n"+
					"#\n"+
					"#   Run `%s prepare db` to initialize the database.\n"+
					"################################################################################\n",
					os.Args[0])
				os.Exit(1)
			}

			if initializingDS, ok := ds.(initializer); ok {
				if err := initializingDS.Initialize(); err != nil {
					initFatal(err, "loading built in data")
				}
			}

			licenseService := license.NewChecker(
				ds,
				"https://kolide.co/api/v0/licenses",
				license.Logger(logger),
			)

			err = licenseService.Start()
			if err != nil {
				initFatal(err, "initializing license service")
			}

			var resultStore kolide.QueryResultStore
			redisPool := pubsub.NewRedisPool(config.Redis.Address, config.Redis.Password)
			resultStore = pubsub.NewRedisQueryResults(redisPool)

			svc, err := service.NewService(ds, resultStore, logger, config, mailService, clock.C, licenseService)
			if err != nil {
				initFatal(err, "initializing service")
			}

			go func() {
				ticker := time.NewTicker(1 * time.Hour)
				for {
					ds.CleanupDistributedQueryCampaigns(time.Now())
					<-ticker.C
				}
			}()

			fieldKeys := []string{"method", "error"}
			requestCount := kitprometheus.NewCounterFrom(prometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys)
			requestLatency := kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys)

			svcLogger := kitlog.With(logger, "component", "service")
			svc = service.NewLoggingService(svc, svcLogger)
			svc = service.NewMetricsService(svc, requestCount, requestLatency)

			httpLogger := kitlog.With(logger, "component", "http")

			var apiHandler, frontendHandler http.Handler
			{
				frontendHandler = prometheus.InstrumentHandler("get_frontend", service.ServeFrontend(httpLogger))
				apiHandler = service.MakeHandler(svc, config.Auth.JwtKey, httpLogger)

				setupRequired, err := service.RequireSetup(svc)
				if err != nil {
					initFatal(err, "fetching setup requirement")
				}
				// WithSetup will check if first time setup is required
				// By performing the same check inside main, we can make server startups
				// more efficient after the first startup.
				if setupRequired {
					apiHandler = service.WithSetup(svc, logger, apiHandler)
					frontendHandler = service.RedirectLoginToSetup(svc, logger, frontendHandler)
				} else {
					frontendHandler = service.RedirectSetupToLogin(svc, logger, frontendHandler)
				}

			}

			// a list of dependencies which could affect the status of the app if unavailable
			healthCheckers := map[string]interface{}{
				"datastore":          ds,
				"query_result_store": resultStore,
			}

			r := http.NewServeMux()
			r.Handle("/healthz", prometheus.InstrumentHandler("healthz", healthz(httpLogger, healthCheckers)))
			r.Handle("/version", prometheus.InstrumentHandler("version", version.Handler()))
			r.Handle("/assets/", prometheus.InstrumentHandler("static_assets", service.ServeStaticAssets("/assets/")))
			r.Handle("/metrics", prometheus.InstrumentHandler("metrics", promhttp.Handler()))
			r.Handle("/api/", apiHandler)
			r.Handle("/", frontendHandler)

			if debug {
				// Add debug endpoints with a random
				// authorization token
				debugToken, err := kolide.RandomText(24)
				if err != nil {
					initFatal(err, "generating debug token")
				}
				r.Handle("/debug/", http.StripPrefix("/debug/", netbug.AuthHandler(debugToken)))
				fmt.Printf("*** Debug mode enabled ***\nAccess the debug endpoints at /debug/?token=%s\n", url.QueryEscape(debugToken))
			}

			srv := &http.Server{
				Addr:              config.Server.Address,
				Handler:           r,
				ReadTimeout:       25 * time.Second,
				WriteTimeout:      40 * time.Second,
				ReadHeaderTimeout: 5 * time.Second,
				IdleTimeout:       5 * time.Minute,
				MaxHeaderBytes:    1 << 18, // 0.25 MB (262144 bytes)
			}
			errs := make(chan error, 2)
			go func() {
				if !config.Server.TLS {
					logger.Log("transport", "http", "address", config.Server.Address, "msg", "listening")
					errs <- srv.ListenAndServe()
				} else {
					logger.Log("transport", "https", "address", config.Server.Address, "msg", "listening")
					srv.TLSConfig = getTLSConfig(config.Server.TLSProfile)
					errs <- srv.ListenAndServeTLS(
						config.Server.Cert,
						config.Server.Key,
					)
				}
			}()
			go func() {
				sig := make(chan os.Signal)
				signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
				<-sig //block on signal
				ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
				errs <- srv.Shutdown(ctx)
			}()

			logger.Log("terminated", <-errs)
			licenseService.Stop()
		},
	}

	serveCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug endpoints")

	return serveCmd
}

// healthz is an http handler which responds with either
// 200 OK if the server can successfuly communicate with it's backends or
// 500 if any of the backends are reporting an issue.
func healthz(logger kitlog.Logger, deps map[string]interface{}) http.HandlerFunc {
	type healthChecker interface {
		HealthCheck() error
	}

	healthy := true
	return func(w http.ResponseWriter, r *http.Request) {
		for name, dep := range deps {
			if hc, ok := dep.(healthChecker); ok {
				err := hc.HealthCheck()
				if err != nil {
					logger.Log("err", err, "health-checker", name)
					healthy = false
				}
			}
		}

		if !healthy {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// Support for TLS security profiles, we set up the TLS configuation based on
// value supplied to server_tls_compatibility command line flag. The default
// profile is 'modern'.
// See https://wiki.mozilla.org/Security/Server_Side_TLS
func getTLSConfig(profile string) *tls.Config {
	cfg := tls.Config{PreferServerCipherSuites: true}

	switch profile {
	case config.TLSProfileModern:
		cfg.MinVersion = tls.VersionTLS12
		cfg.CurvePreferences = append(cfg.CurvePreferences,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		)
		cfg.CipherSuites = append(cfg.CipherSuites,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		)
	case config.TLSProfileIntermediate:
		cfg.MinVersion = tls.VersionTLS10
		cfg.CurvePreferences = append(cfg.CurvePreferences,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		)
		cfg.CipherSuites = append(cfg.CipherSuites,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		)
	case config.TLSProfileOld:
		cfg.MinVersion = tls.VersionSSL30
		cfg.CurvePreferences = append(cfg.CurvePreferences,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		)
		cfg.CipherSuites = append(cfg.CipherSuites,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		)
	default:
		panic("invalid tls profile " + profile)
	}

	return &cfg
}
