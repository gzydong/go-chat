package provider

import "go-chat/internal/pkg/email"

type Providers struct {
	EmailClient *email.Client
}

func NewProviders(emailClient *email.Client) *Providers {
	return &Providers{EmailClient: emailClient}
}
