package config

// OAuth 第三方登录配置
type OAuth struct {
	Github *GithubOAuth `json:"github" yaml:"github"`
	Gitee  *GiteeOAuth  `json:"gitee" yaml:"gitee"`
}

// GithubOAuth GitHub OAuth配置
type GithubOAuth struct {
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	RedirectURL  string `json:"redirect_uri" yaml:"redirect_uri"`
}

// GiteeOAuth Gitee OAuth配置
type GiteeOAuth struct {
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	RedirectURL  string `json:"redirect_uri" yaml:"redirect_uri"`
}
