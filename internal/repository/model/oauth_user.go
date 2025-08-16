package model

import "time"

// OAuthType 第三方登录类型
type OAuthType string

const (
	OAuthTypeGithub OAuthType = "github"
	OAuthTypeGitee  OAuthType = "gitee"
)

// OAuthUser 第三方用户信息
type OAuthUser struct {
	Id            int32     `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"` // 自增ID
	UserId        int32     `json:"user_id" gorm:"column:user_id;"`                 // 用户ID
	OAuthType     OAuthType `json:"oauth_type" gorm:"column:oauth_type;"`           // 第三方类型
	OAuthId       string    `json:"oauth_id" gorm:"column:oauth_id;"`               // 第三方用户ID
	Username      string    `json:"username" gorm:"column:username;"`               // 第三方用户名
	Nickname      string    `json:"nickname" gorm:"column:nickname;"`               // 第三方昵称
	Email         string    `json:"email" gorm:"column:email;"`                     // 邮箱
	Avatar        string    `json:"avatar" gorm:"column:avatar;"`                   // 头像
	AccessToken   string    `json:"access_token" gorm:"column:access_token;"`       // 访问令牌
	RefreshToken  string    `json:"refresh_token" gorm:"column:refresh_token;"`     // 刷新令牌
	ExpiresIn     int64     `json:"expires_in" gorm:"column:expires_in;"`           // 过期时间
	Scope         string    `json:"scope" gorm:"column:scope;"`                     // 授权范围
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;"`           // 创建时间
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;"`           // 更新时间
	LastLoginTime time.Time `json:"last_login_time" gorm:"column:last_login_time;"` // 最后登录时间
}

// TableName 表名
func (OAuthUser) TableName() string {
	return "oauth_users"
}
