package provider

import (
	"context"
	"net"
	"net/http"
	"time"
)

const timeout = 20 * time.Second

func dialTimeout(ctx context.Context, network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialTLSContext:        dialTimeout,
			ResponseHeaderTimeout: time.Second * 2,
		},
		Timeout: timeout,
	}
}
