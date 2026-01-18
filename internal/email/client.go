package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type Client struct {
	host     string
	port     string
	username string
	password string
	from     string
	auth     smtp.Auth
}

func NewMailtrapClient() *Client {
	host := os.Getenv("MAILTRAP_HOST")
	port := os.Getenv("MAILTRAP_PORT")
	username := os.Getenv("MAILTRAP_USER")
	password := os.Getenv("MAILTRAP_PASS")
	from := os.Getenv("MAILTRAP_FROM")

	auth := smtp.PlainAuth("", username, password, host)

	return &Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		auth:     auth,
	}
}

func (c *Client) SendMail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", c.host, c.port)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		c.from, to, subject, body,
	)

	return smtp.SendMail(addr, c.auth, c.from, []string{to}, []byte(msg))
}
