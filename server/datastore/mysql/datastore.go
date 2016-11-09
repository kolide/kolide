package mysql

import (
	"fmt"
	"strings"

	"github.com/WatchBeam/clock"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql" // db driver
	"github.com/jmoiron/sqlx"
	"github.com/kolide/kolide-ose/server/kolide"
)

const (
	defaultSelectLimit = 1000
)

// Datastore is an implementation of kolide.Datastore interface backed by
// MySQL
type Datastore struct {
	db     *sqlx.DB
	logger log.Logger
	clock  clock.Clock
}

// New creates an MySQL datastore.
func New(dbConnectString string, c clock.Clock, opts ...DBOption) (*Datastore, error) {
	var (
		ds  *Datastore
		err error
		db  *sqlx.DB
	)

	options := dbOptions{
		maxAttempts: defaultMaxAttempts,
		logger:      log.NewNopLogger(),
	}

	for _, setOpt := range opts {
		setOpt(&options)
	}

	for attempt := 0; attempt < options.maxAttempts; attempt++ {
		if db, err = sqlx.Connect("mysql", dbConnectString); err == nil {
			break
		}
	}

	if db == nil {
		return nil, err
	}

	ds = &Datastore{db, options.logger, c}

	return ds, nil

}

func (d *Datastore) Name() string {
	return "mysql"
}

// Migrate creates database
func (d *Datastore) Migrate() error {
	var (
		err error
		sql []byte
	)

	if sql, err = Asset("db/up.sql"); err != nil {
		return err
	}

	tx := d.db.MustBegin()

	for _, statement := range strings.SplitAfter(string(sql), ";") {
		if _, err = tx.Exec(statement); err != nil {
			if err.Error() != "Error 1065: Query was empty" {
				tx.Rollback()
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

// Drop removes database
func (d *Datastore) Drop() error {
	var (
		sql []byte
		err error
	)

	if sql, err = Asset("db/down.sql"); err != nil {
		return err
	}

	tx := d.db.MustBegin()

	for _, statement := range strings.SplitAfter(string(sql), ";") {
		if _, err = tx.Exec(statement); err != nil {
			if err.Error() != "Error 1065: Query was empty" {
				tx.Rollback()
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

// Close frees resources associated with underlying mysql connection
func (d *Datastore) Close() error {
	return d.db.Close()
}

func (d *Datastore) log(msg string) {
	d.logger.Log("comp", d.Name(), "msg", msg)
}

func appendListOptionsToSQL(sql string, opts kolide.ListOptions) string {
	if opts.OrderKey != "" {
		direction := "ASC"
		if opts.OrderDirection == kolide.OrderDescending {
			direction = "DESC"
		}

		sql = fmt.Sprintf("%s ORDER BY %s %s", sql, opts.OrderKey, direction)
	}
	// REVIEW: If caller doesn't supply a limit apply a default limit of 1000
	// to insure that an unbounded query with many results doesn't consume too
	// much memory or hang
	if opts.PerPage == 0 {
		opts.PerPage = defaultSelectLimit
	}

	sql = fmt.Sprintf("%s LIMIT %d", sql, opts.PerPage)

	offset := opts.PerPage * opts.Page

	if offset > 0 {
		sql = fmt.Sprintf("%s OFFSET %d", sql, offset)
	}

	return sql
}
