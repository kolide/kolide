package kolide

import (
	"bytes"
	"html/template"
	"time"
)

// PasswordResetStore manages password resets in the Datastore
type PasswordResetStore interface {
	NewPasswordResetRequest(req *PasswordResetRequest) (*PasswordResetRequest, error)
	SavePasswordResetRequest(req *PasswordResetRequest) error
	DeletePasswordResetRequest(req *PasswordResetRequest) error
	DeletePasswordResetRequestsForUser(userID uint) error
	FindPassswordResetByID(id uint) (*PasswordResetRequest, error)
	FindPassswordResetsByUserID(id uint) ([]*PasswordResetRequest, error)
	FindPassswordResetByToken(token string) (*PasswordResetRequest, error)
	FindPassswordResetByTokenAndUserID(token string, id uint) (*PasswordResetRequest, error)
}

// Mailer is an email campaign
// Types which implement the Campaign interface
// can be marshalled into an email body
type Mailer interface {
	Message() ([]byte, error)
}

type Email struct {
	Subject string
	To      []string
	Config  *SMTPConfig
	Mailer  Mailer
}

type MailService interface {
	SendEmail(e Email) error
}

// PasswordResetRequest represents a database table for
// Password Reset Requests
type PasswordResetRequest struct {
	UpdateCreateTimestamps
	ID        uint
	ExpiresAt time.Time `db:"expires_at"`
	UserID    uint      `db:"user_id"`
	Token     string
}

const passwordResetTemplate = `
You requested a password reset,
Follow the link below to reset your password:
http://localhost:8080/login/reset?token={{.Token}}
`

type PasswordResetMailer struct {
	// URL for the Kolide application
	KolideServerURL string
	// Token password reset token
	Token string
}

func (r PasswordResetMailer) Message() ([]byte, error) {
	t, err := getTemplate("server/mail/templates/password_reset.html")
	if err != nil {
		return nil, err
	}

	var msg bytes.Buffer
	if err = t.Execute(&msg, r); err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}

func getTemplate(templatePath string) (*template.Template, error) {
	templateData, err := Asset(templatePath)
	if err != nil {
		return nil, err
	}

	t, err := template.New("email_template").Parse(string(templateData))
	if err != nil {
		return nil, err
	}

	return t, nil
}
