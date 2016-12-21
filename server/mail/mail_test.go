package mail

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
)

type mockMailer struct{}

func (m *mockMailer) SendEmail(e kolide.Email) error {
	return nil
}

func getMailer() kolide.MailService {

	if os.Getenv("MAIL_TEST") == "" {
		return &mockMailer{}
	}
	return NewService()
}

func functionName(f func(*testing.T, kolide.MailService)) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	elements := strings.Split(fullName, ".")
	return elements[len(elements)-1]
}

var testFunctions = [...]func(*testing.T, kolide.MailService){
	testSMTPPlainAuth,
	testSMTPSkipVerify,
	testSMTPNoAuth,
}

func TestMail(t *testing.T) {
	for _, f := range testFunctions {
		r := getMailer()

		t.Run(functionName(f), func(t *testing.T) {
			f(t, r)
		})
	}
}

func testSMTPPlainAuth(t *testing.T, mailer kolide.MailService) {
	mail := kolide.Email{
		Subject: "smtp plain auth",
		To:      []string{"john@kolide.co"},
		Config: &kolide.AppConfig{
			SMTPConfigured:           true,
			SMTPDisabled:             false,
			SMTPAuthenticationType:   kolide.AuthTypeUserNamePassword,
			SMTPAuthenticationMethod: kolide.AuthMethodPlain,
			SMTPUserName:             "bob",
			SMTPPassword:             "secret",
			SMTPEnableTLS:            true,
			SMTPVerifySSLCerts:       true,
			SMTPEnableStartTLS:       true,
			SMTPPort:                 1025,
			SMTPServer:               "localhost",
			SMTPSenderAddress:        "kolide@kolide.com",
		},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: "https://localhost:8080",
		},
	}

	err := mailer.SendEmail(mail)
	assert.Nil(t, err)
}

func testSMTPSkipVerify(t *testing.T, mailer kolide.MailService) {
	mail := kolide.Email{
		Subject: "skip verify",
		To:      []string{"john@kolide.co"},
		Config: &kolide.AppConfig{
			SMTPConfigured:           true,
			SMTPDisabled:             false,
			SMTPAuthenticationType:   kolide.AuthTypeUserNamePassword,
			SMTPAuthenticationMethod: kolide.AuthMethodPlain,
			SMTPUserName:             "bob",
			SMTPPassword:             "secret",
			SMTPEnableTLS:            true,
			SMTPVerifySSLCerts:       false,
			SMTPEnableStartTLS:       true,
			SMTPPort:                 1025,
			SMTPServer:               "localhost",
			SMTPSenderAddress:        "kolide@kolide.com",
		},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: "https://localhost:8080",
		},
	}

	err := mailer.SendEmail(mail)
	assert.Nil(t, err)
}

func testSMTPNoAuth(t *testing.T, mailer kolide.MailService) {
	mail := kolide.Email{
		Subject: "no auth",
		To:      []string{"bob@foo.com"},
		Config: &kolide.AppConfig{
			SMTPConfigured:         true,
			SMTPDisabled:           false,
			SMTPAuthenticationType: kolide.AuthTypeNone,
			SMTPEnableTLS:          true,
			SMTPVerifySSLCerts:     true,
			SMTPPort:               1025,
			SMTPServer:             "localhost",
			SMTPSenderAddress:      "kolide@kolide.com",
		},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: "https://localhost:8080",
		},
	}

	err := mailer.SendEmail(mail)
	assert.Nil(t, err)
}
