package email

import (
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type Client struct {
	api *brevo.APIClient
}

func NewClient() *Client {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", os.Getenv("BREVO_API_KEY"))

	return &Client{
		api: brevo.NewAPIClient(cfg),
	}
}
