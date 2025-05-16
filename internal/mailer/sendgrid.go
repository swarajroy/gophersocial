package mailer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail       string
	apiKey          string
	client          *sendgrid.Client
	templateBuilder *TemplateBuilder
}

func NewSendGridMailer(fromEmail, apiKey string, templateBuilder *TemplateBuilder) (*SendGridMailer, error) {
	if apiKey == "" {
		return nil, ErrApiKeyRequired
	}
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail:       fromEmail,
		apiKey:          apiKey,
		client:          client,
		templateBuilder: templateBuilder,
	}, nil
}

func (m *SendGridMailer) Send(ctx context.Context, templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	//template parsing and building
	subject, body, err := getSubjectAndBody(templateFile, data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := range maxRetires {
		resp, err := m.client.SendWithContext(ctx, message)
		if err != nil {
			log.Printf("failed to send email to %+v, %d attempt of %d\n", email, i+1, maxRetires)
			log.Printf("Error: %+v\n", err.Error())
			//exponential back off
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email %+v sent with status code %d\n", email, resp.StatusCode)
		return 200, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempts", maxRetires)
}
