package provider

import "go-chat/internal/pkg/email"

type Providers struct {
	EmailClient *email.Client
}
