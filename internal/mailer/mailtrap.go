package mailer

import (
	"context"
	"fmt"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

type MailtrapMailer struct {
	fromEmail       string
	apiKey          string
	templateBuilder *TemplateBuilder
}

func NewMailtrapMailer(fromEmail, apiKey string, templateBuilder *TemplateBuilder) (*MailtrapMailer, error) {
	if apiKey == "" {
		return nil, ErrApiKeyRequired
	}

	return &MailtrapMailer{
		fromEmail:       fromEmail,
		apiKey:          apiKey,
		templateBuilder: templateBuilder,
	}, nil
}

func (mt *MailtrapMailer) Send(ctx context.Context, templateFile, username, email string, data any, isSandbox bool) (int, error) {
	//template parsing and building
	res, err := mt.templateBuilder.build(templateFile, data)
	if err != nil {
		return -1, err
	}

	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", mt.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", res.subject.String())

	// Set the HTML version of the email
	message.AddAlternative("text/html", res.body.String())

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", mt.apiKey)

	// Send the email with retries
	for i := range maxRetires {
		if !isSandbox {
			if err := dialer.DialAndSend(message); err != nil {

				log.Printf("failed to send email to %+v, %d attempt of %d\n", email, i+1, maxRetires)
				log.Printf("Error: %+v\n", err.Error())
				//exponential back off
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			log.Printf("Email %+v sent with status code %d\n", email, 200)
			return 200, nil
		}
	}

	return -1, fmt.Errorf("failed to send email after %d attempts", maxRetires)
}
