package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type IClient interface {
	GetAuthUserInfo(ctx context.Context, accessToken string) (*AuthUserInfo, error)
}

type Client struct {
	c *http.Client
}

func NewClient(c *http.Client) IClient {
	return &Client{
		c: c,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("status: %s, message: %s", e.Status, e.Message)
}

// GetAuthUserInfo 获取授权用户的信息
func (g *Client) GetAuthUserInfo(ctx context.Context, accessToken string) (*AuthUserInfo, error) {
	return call[AuthUserInfo](ctx, "/user", "GET", accessToken)
}

func call[T any](ctx context.Context, uri string, method string, accessToken string) (*T, error) {
	url := "https://api.github.com" + uri

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}

	if errResp.Status != "" {
		return nil, errResp
	}

	var resp T
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type AuthUserInfo struct {
	Login             string      `json:"login"`
	Id                int         `json:"id"`
	NodeId            string      `json:"node_id"`
	AvatarUrl         string      `json:"avatar_url"`
	GravatarId        string      `json:"gravatar_id"`
	Url               string      `json:"url"`
	HtmlUrl           string      `json:"html_url"`
	FollowersUrl      string      `json:"followers_url"`
	FollowingUrl      string      `json:"following_url"`
	GistsUrl          string      `json:"gists_url"`
	StarredUrl        string      `json:"starred_url"`
	SubscriptionsUrl  string      `json:"subscriptions_url"`
	OrganizationsUrl  string      `json:"organizations_url"`
	ReposUrl          string      `json:"repos_url"`
	EventsUrl         string      `json:"events_url"`
	ReceivedEventsUrl string      `json:"received_events_url"`
	Type              string      `json:"type"`
	UserViewType      string      `json:"user_view_type"`
	SiteAdmin         bool        `json:"site_admin"`
	Name              string      `json:"name"`
	Company           interface{} `json:"company"`
	Blog              string      `json:"blog"`
	Location          interface{} `json:"location"`
	Email             *string     `json:"email"`
	Hireable          interface{} `json:"hireable"`
	Bio               interface{} `json:"bio"`
	TwitterUsername   interface{} `json:"twitter_username"`
	NotificationEmail interface{} `json:"notification_email"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}
