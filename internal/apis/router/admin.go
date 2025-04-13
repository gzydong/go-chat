package router

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"go-chat/internal/apis/handler/admin"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/jwtutil"
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
	)

	// v1 接口
	v1 := router.Group("/admin/v1")
	{
		index := v1.Group("/index")
		{
			index.GET("", core.HandlerFunc(handler.V1.Index.Index))
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", core.HandlerFunc(handler.V1.Auth.Login))
			auth.GET("/captcha", core.HandlerFunc(handler.V1.Auth.Captcha))
			auth.GET("/logout", authorize, core.HandlerFunc(handler.V1.Auth.Logout))
		}
	}
}
