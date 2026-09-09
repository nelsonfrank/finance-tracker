package mailer

import (
	"bytes"
	"errors"
	"html/template"

	"github.com/resend/resend-go/v3"
)

type ResendClient struct {
	fromEmail string
	apiKey    string
}

func NewResendClient(apiKey, fromEmail string) (ResendClient, error) {
	if apiKey == "" {
		return ResendClient{}, errors.New("api key is required")
	}

	if fromEmail == "" {
		return ResendClient{}, errors.New("from email is required")
	}

	return ResendClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m ResendClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	client := resend.NewClient(m.apiKey)

	params := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{email},
		Html:    body.String(),
		Subject: subject.String(),
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		return -1, err
	}

	return 200, nil
}
