package email

type Sender struct {
	client *Client
}

func NewSender(client *Client) *Sender {
	return &Sender{client: client}
}

func (s *Sender) Send(to, subject, body string) error {
	return s.client.SendMail(to, subject, body)
}

func (s *Sender) SendResetPassword(to, resetURL string) error {
	body := ResetPasswordTemplate(resetURL)
	return s.Send(to, "Reset your password", body)
}
