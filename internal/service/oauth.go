package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gzydong/go-chat/config"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/thirdsdk/gitee"
	"github.com/gzydong/go-chat/internal/pkg/thirdsdk/github"
	"github.com/gzydong/go-chat/internal/pkg/utils"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"

	"golang.org/x/oauth2"
)

// IOAuthService OAuth服务接口
type IOAuthService interface {
	// GetAuthURL 获取第三方登录授权URL
	GetAuthURL(ctx context.Context, oauthType model.OAuthType) (string, error)

	// HandleCallback 处理第三方登录回调
	HandleCallback(ctx context.Context, oauthType model.OAuthType, code, state string) (*model.OAuthUser, error)
}

// OAuthService 第三方登录服务
type OAuthService struct {
	OauthUsers   *repo.OAuthUsers
	UserSrv      *UserService
	Config       *config.Config
	GiteeClient  gitee.IClient
	GithubClient github.IClient
	Redis        *redis.Client
}

// GetAuthURL 获取第三方登录授权URL
func (s *OAuthService) GetAuthURL(ctx context.Context, oauthType model.OAuthType) (string, error) {
	var conf *oauth2.Config

	switch oauthType {
	case model.OAuthTypeGithub:
		conf = s.getGithubConfig()
	case model.OAuthTypeGitee:
		conf = s.getGiteeConfig()
	default:
		return "", entity.ErrOauthTypeInvalid
	}

	state := strings.ReplaceAll(uuid.New().String(), "-", "")

	err := s.Redis.SetEx(ctx, fmt.Sprintf("oauth:state:%s", state), "1", time.Minute).Err()
	if err != nil {
		return "", err
	}

	return conf.AuthCodeURL(state), nil
}

// HandleCallback 处理第三方登录回调
func (s *OAuthService) HandleCallback(ctx context.Context, oauthType model.OAuthType, code, state string) (*model.OAuthUser, error) {
	if s.Redis.Get(ctx, fmt.Sprintf("oauth:state:%s", state)).Val() != "1" {
		return nil, entity.ErrStateInvalid
	}

	var conf *oauth2.Config
	var oauthUser *model.OAuthUser
	var err error

	switch oauthType {
	case model.OAuthTypeGithub:
		conf = s.getGithubConfig()
		oauthUser, err = s.handleGithubCallback(ctx, conf, code)
	case model.OAuthTypeGitee:
		conf = s.getGiteeConfig()
		oauthUser, err = s.handleGiteeCallback(ctx, conf, code)
	default:
		return nil, entity.ErrOauthTypeInvalid
	}

	if err != nil {
		return nil, err
	}

	// 查找是否已存在该第三方账号绑定
	existUser, err := s.OauthUsers.FindByTypeAndId(ctx, oauthType, oauthUser.OAuthId)
	if err != nil && !utils.IsSqlNoRows(err) {
		return nil, err
	}

	if existUser != nil {
		// 更新令牌信息
		updates := map[string]any{
			"access_token":    oauthUser.AccessToken,
			"refresh_token":   oauthUser.RefreshToken,
			"expires_in":      oauthUser.ExpiresIn,
			"last_login_time": time.Now(),
		}

		_, err = s.OauthUsers.UpdateById(ctx, existUser.Id, updates)
		if err != nil {
			return nil, err
		}
		return existUser, nil
	} else {
		oauthUser.LastLoginTime = time.Now()
		// 新用户，需要创建
		if err := s.OauthUsers.Create(ctx, oauthUser); err != nil {
			return nil, err
		}
	}

	// 新用户，需要创建
	return oauthUser, nil
}

// getGithubConfig 获取GitHub OAuth配置
func (s *OAuthService) getGithubConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.Config.OAuth.Github.ClientID,
		ClientSecret: s.Config.OAuth.Github.ClientSecret,
		RedirectURL:  s.Config.OAuth.Github.RedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

// getGiteeConfig 获取Gitee OAuth配置
func (s *OAuthService) getGiteeConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.Config.OAuth.Gitee.ClientID,
		ClientSecret: s.Config.OAuth.Gitee.ClientSecret,
		RedirectURL:  s.Config.OAuth.Gitee.RedirectURL,
		Scopes:       []string{"user_info", "emails"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://gitee.com/oauth/authorize",
			TokenURL: "https://gitee.com/oauth/token",
		},
	}
}

// handleGithubCallback 处理GitHub回调
func (s *OAuthService) handleGithubCallback(ctx context.Context, conf *oauth2.Config, code string) (*model.OAuthUser, error) {
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("获取GitHub token失败: %w", err)
	}

	userInfo, err := s.GithubClient.GetAuthUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("获取GitHub用户信息失败: %w", err)
	}

	oauthUser := &model.OAuthUser{
		OAuthType:    model.OAuthTypeGithub,
		OAuthId:      fmt.Sprintf("%d", userInfo.Id),
		Username:     userInfo.Login,
		Nickname:     userInfo.Name,
		Email:        lo.FromPtr(userInfo.Email),
		Avatar:       userInfo.AvatarUrl,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry.Unix(),
		Scope:        "user:email",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 如果没有昵称，使用用户名
	if oauthUser.Nickname == "" {
		oauthUser.Nickname = oauthUser.Username
	}

	return oauthUser, nil
}

// handleGiteeCallback 处理Gitee回调
func (s *OAuthService) handleGiteeCallback(ctx context.Context, conf *oauth2.Config, code string) (*model.OAuthUser, error) {
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("获取Gitee token失败: %w", err)
	}

	userInfo, err := s.GiteeClient.GetAuthUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("获取Gitee用户信息失败: %w", err)
	}

	oauthUser := &model.OAuthUser{
		OAuthType:    model.OAuthTypeGitee,
		OAuthId:      strconv.Itoa(userInfo.Id),  // 应从API响应中获取
		Username:     userInfo.Login,             // 应从API响应中获取
		Nickname:     userInfo.Name,              // 应从API响应中获取
		Email:        lo.FromPtr(userInfo.Email), // 应从API响应中获取
		Avatar:       userInfo.AvatarUrl,         // 应从API响应中获取
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry.Unix(),
		Scope:        "user_info emails",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return oauthUser, nil
}
