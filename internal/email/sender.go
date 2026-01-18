package email

import (
	"context"
	"fmt"
	"log"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type Sender struct {
	client *Client
	from   string
}

func NewSender(client *Client) *Sender {
	return &Sender{
		client: client,
		from:   os.Getenv("SENDER_EMAIL"),
	}
}

func (s *Sender) Send(
	to string,
	subject string,
	htmlContent string,
) error {
	email := brevo.SendSmtpEmail{
		To: []brevo.SendSmtpEmailTo{
			{Email: to},
		},
		Subject:     subject,
		HtmlContent: htmlContent,
		Sender: &brevo.SendSmtpEmailSender{
			Email: s.from,
			Name:  "Auth Service",
		},
	}

	_, _, err := s.client.api.TransactionalEmailsApi.
		SendTransacEmail(context.Background(), email)

	if err != nil {
		return fmt.Errorf("send email failed: %w", err)
	}

	return nil
}

func (s *Sender) SendResetPassword(to, resetURL string) error {
	log.Print(s.from)
	html := ResetPasswordTemplate(resetURL)
	return s.Send(to, "Reset your password", html)
}
