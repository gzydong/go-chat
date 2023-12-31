package provider

import (
	"net/http"
	"time"

	"go-chat/internal/pkg/client"
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func NewRequestClient(c *http.Client) *client.RequestClient {
	return client.NewRequestClient(c)
}
