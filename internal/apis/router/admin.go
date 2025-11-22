package router

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	admin2 "github.com/gzydong/go-chat/api/pb/admin/v1"
	"github.com/gzydong/go-chat/internal/apis/handler/admin"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/jwtutil"
)

// RegisterAdminRoute 注册 Admin 路由
func RegisterAdminRoute(secret string, router *gin.Engine, handler *admin.Handler, storage middleware.IStorage) {
	// 授权验证中间件
	authorize := middleware.NewJwtMiddleware[entity.AdminClaims](
		[]byte(secret), storage,
		func(ctx context.Context, claims *jwtutil.JwtClaims[entity.AdminClaims]) error {
			if claims.RegisteredClaims.Issuer != entity.JwtIssuerAdmin {
				return errors.New("授权异常，请登录后操作")
			}

			adminInfo, err := handler.AdminRepo.FindById(ctx, claims.Metadata.AdminId)
			if err != nil {
				return err
			}

			if adminInfo.IsDisabled() {
				return entity.ErrAccountDisabled
			}

			return nil
		},
		func(option *middleware.JwtMiddlewareOption) {
			option.ExclusionPaths = []string{
				"/backend/auth/login",
				"/backend/auth/captcha",
			}
		},
	)

	resp := &Interceptor{}

	group := router.Group("/").Use(authorize)

	admin2.RegisterTotpHandler(group, resp, handler.Totp)
	admin2.RegisterAuthHandler(group, resp, handler.Auth)
	admin2.RegisterMenuHandler(group, resp, handler.Menu)
	admin2.RegisterResourceHandler(group, resp, handler.Resource)
	admin2.RegisterAdminHandler(group, resp, handler.Admin)
	admin2.RegisterRoleHandler(group, resp, handler.Role)
	admin2.RegisterUserHandler(group, resp, handler.User)

	registerAdminCustomApiRouter(resp, group, handler)

}

func registerAdminCustomApiRouter(resp *Interceptor, api gin.IRoutes, handler *admin.Handler) {
	api.POST("/backend/upload/file", HandlerFunc(resp, func(c *gin.Context) (any, error) {
		return map[string]any{
			"url": "https://www.cox.com/" + uuid.NewString(),
		}, nil
	}))
}
