package mailer

import (
	"bytes"
	"text/template"
)

type TemplateBuilder struct {
}

func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{}
}

type TemplateBuilderResultContainer struct {
	subject, body *bytes.Buffer
}

func (tb *TemplateBuilder) build(templateFile string, data any) (*TemplateBuilderResultContainer, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return nil, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return nil, err
	}

	return &TemplateBuilderResultContainer{
		subject: subject,
		body:    body,
	}, nil
}
