package mailer

import "embed"

const (
	FromName               = "Gosocial"
	MaxRetries             = 3
	AccountWelcomeTemplate = "accountConfirmation.tmpl"
)

var FS embed.FS

type Client interface {
	Send(templateFile string, username, email string, data any, isSanbox bool) error
}
