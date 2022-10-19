package provider

import (
	"net/http"
	"time"

	"go-chat/internal/pkg/client"
)

const timeout = 5 * time.Second

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func NewRequestClient(c *http.Client) *client.RequestClient {
	return client.NewRequestClient(c)
}
