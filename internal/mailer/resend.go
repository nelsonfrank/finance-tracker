package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/resend/resend-go/v2"
)

type resendClient struct {
	fromEmail string
	apiKey    string
	client *resend.Client
}

func NewResendClient(apiKey, fromEmail string) (resendClient, error) {
	if apiKey == "" {
		return resendClient{}, errors.New("api key is required")
	}

	client := resend.NewClient(apiKey)
	return resendClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client: client,
	}, nil
}


func (m resendClient) Send(templateFile, username, email string, data any, isSandbox bool)(int, error) {
 

  // template parsing and building
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

  params := &resend.SendEmailRequest{
      From:        "Financing App <onboarding@resend.dev>",
      To:          []string{username},
      Subject:     subject.String(),
      Html:        body.String(),
  }

  sent, err := m.client.Emails.Send(params)

  if err != nil {
    return -1, err
  }
  fmt.Println(sent.Id)
  return 200, nil
}