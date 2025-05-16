package mailer

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"text/template"
)

const (
	FromName               = "GopherSocial"
	maxRetires             = 3
	UserInvitationTemplate = "user_invitation.tmpl"
)

var (
	ErrApiKeyRequired = errors.New("apiKey is required")
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(ctx context.Context, templateFile, username, email string, data any, isSandbox bool) (int, error)
}

func getSubjectAndBody(templateFile string, data any) (*bytes.Buffer, *bytes.Buffer, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return nil, nil, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, nil, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return nil, nil, err
	}

	return subject, body, nil
}
