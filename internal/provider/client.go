package provider

import (
	"net/http"
	"time"

	"go-chat/internal/pkg/ipaddress"
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func NewIpAddressClient(c *http.Client) *ipaddress.Client {
	return ipaddress.NewClient(c)
}
