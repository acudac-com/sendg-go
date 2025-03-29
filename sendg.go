package sendgrid

import (
	"context"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.alis.build/alog"
)

type Client struct {
	client *sendgrid.Client
}

type Template[T any] struct {
	id     string
	client *Client
}

func NewClient() *Client {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		alog.Fatal(context.Background(), "Missing SENDGRID_API_KEY environment variable")
	}
	return &Client{
		sendgrid.NewSendClient(apiKey),
	}
}

func NewTemplate[T any](client *Client, id string, emptyData T) *Template[T] {
	return &Template[T]{
		id, client,
	}
}

// Sends an email via Sendgrid, using the template.
//   - tos: mulitple addresses to send the emails to
//   - from: the email from which to send the email. This needs to be a verified email within your sendgrid account.  More details on authenticating
//     your domain is available at: https://app.sendgrid.com/settings/sender_auth
func (t *Template[T]) Mail(from string, data T, tos ...string) error {
	personalization := mail.NewPersonalization()
	personalization.To = []*mail.Email{}
	for _, to := range tos {
		personalization.To = append(personalization.To, mail.NewEmail("", to))
	}
	personalization.SetDynamicTemplateData("data", data)

	// Create email
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail("", from))
	message.AddPersonalizations(personalization)
	message.SetTemplateID(t.id)

	// Send email
	_, err := t.client.client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
