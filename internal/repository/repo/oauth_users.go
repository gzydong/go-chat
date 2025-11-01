package repo

import (
	"context"

	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type OAuthUsers struct {
	core.Repo[model.OAuthUser]
}

func NewOAuthUsers(db *gorm.DB) *OAuthUsers {
	return &OAuthUsers{Repo: core.NewRepo[model.OAuthUser](db)}
}

// FindByTypeAndId 根据第三方类型和ID查找用户
func (o *OAuthUsers) FindByTypeAndId(ctx context.Context, oauthType model.OAuthType, oauthId string) (*model.OAuthUser, error) {
	return o.FindByWhere(ctx, "oauth_type = ? AND oauth_id = ?", oauthType, oauthId)
}

// FindAllByUserId 根据用户ID查找OAuth绑定信息
func (o *OAuthUsers) FindAllByUserId(ctx context.Context, userId int64) ([]*model.OAuthUser, error) {
	return o.FindAllByWhere(ctx, "user_id = ?", userId)
}
