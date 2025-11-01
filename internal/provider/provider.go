package provider

import "github.com/gzydong/go-chat/internal/pkg/email"

type Providers struct {
	EmailClient *email.Client
}
