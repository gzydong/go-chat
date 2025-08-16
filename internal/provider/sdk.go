package provider

import (
	"net/http"

	"go-chat/internal/pkg/thirdsdk/gitee"
	"go-chat/internal/pkg/thirdsdk/github"
)

func NewGiteeClient(c *http.Client) gitee.IClient {
	return gitee.NewClient(c)
}

func NewGithubClient(c *http.Client) github.IClient {
	return github.NewClient(c)
}
