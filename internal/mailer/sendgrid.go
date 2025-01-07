package mailer

import (
	"bytes"
	"time"
	"fmt"
	"log"
	"text/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any isSandbox bool) error {
	from := mail.NewMail(fromName)
	to := mail.NewEmail(username, emai)

	tmpl, err := template.ParseFs(FS, "templates/"+ AccountWelcomeTemplate)
	if err != nil {
		return err	
	}
	
	subject := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err	
	}

	body := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err	
	}

	message := mail.NewSingleEmail(from, subject, to, "", body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		}
	})

	for i := 0; i < MaxRetries; i++ {
		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d ", email, i+1, MaxRetries)
			log.Printf("Error %v", err.Error())
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("Email sent with status code %v", email, response.StatusCode)
		return nil
	}

	return ftm.Errorf("failed to send email after %d attemps", MaxRetries)
}