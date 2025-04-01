package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile string, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, fmt.Errorf("failed to parse template: %w", err)
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return -1, fmt.Errorf("failed to execute subject template: %w", err)
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return -1, fmt.Errorf("failed to execute body template: %w", err)
	}

	// If either subject or body is empty, log it for debugging
	if subject.Len() == 0 {
		log.Printf("Warning: Email subject is empty for recipient %s", email)
	}
	if body.Len() == 0 {
		log.Printf("Warning: Email body is empty for recipient %s", email)
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var retryErr error
	for i := 0; i < maxRetries; i++ {
		response, retryErr := m.client.Send(message)
		if retryErr != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		if response.StatusCode >= 400 {
			log.Printf("Email send failed with status code %d for %s. Response: %v",
				response.StatusCode, email, response.Body)
			return -1, fmt.Errorf("sendgrid error: status code %d", response.StatusCode)
		}

		return response.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempts,error:%v", maxRetries, retryErr)
}
