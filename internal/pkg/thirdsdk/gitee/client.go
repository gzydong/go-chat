package gitee

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
	Code    int    `json:"code"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

type AuthUserInfo struct {
	Id                int       `json:"id"`
	Login             string    `json:"login"`
	Name              string    `json:"name"`
	AvatarUrl         string    `json:"avatar_url"`
	Url               string    `json:"url"`
	HtmlUrl           string    `json:"html_url"`
	Remark            string    `json:"remark"`
	FollowersUrl      string    `json:"followers_url"`
	FollowingUrl      string    `json:"following_url"`
	GistsUrl          string    `json:"gists_url"`
	StarredUrl        string    `json:"starred_url"`
	SubscriptionsUrl  string    `json:"subscriptions_url"`
	OrganizationsUrl  string    `json:"organizations_url"`
	ReposUrl          string    `json:"repos_url"`
	EventsUrl         string    `json:"events_url"`
	ReceivedEventsUrl string    `json:"received_events_url"`
	Type              string    `json:"type"`
	Bio               string    `json:"bio"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	Stared            int       `json:"stared"`
	Watched           int       `json:"watched"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Email             *string   `json:"email"`
}

// GetAuthUserInfo 获取授权用户的信息
func (g *Client) GetAuthUserInfo(ctx context.Context, accessToken string) (*AuthUserInfo, error) {
	return call[AuthUserInfo](ctx, fmt.Sprintf("/api/v5/user?access_token=%s", accessToken), "GET", nil)
}

func call[T any](ctx context.Context, uri string, method string, req any) (*T, error) {
	url := "https://gitee.com" + uri

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	// nolint
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}

	if errResp.Code != 0 {
		return nil, errResp
	}

	if res.StatusCode != http.StatusOK {
		errResp.Code = res.StatusCode
		return nil, errResp
	}

	var resp T
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
