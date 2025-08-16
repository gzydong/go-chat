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
		auth := v1.Group("/auth")
		{
			auth.POST("/login", core.HandlerFunc(handler.Auth.Login))
			auth.POST("/captcha", core.HandlerFunc(handler.Auth.Captcha))
			auth.GET("/logout", authorize, core.HandlerFunc(handler.Auth.Logout))
			auth.POST("/detail", authorize, core.HandlerFunc(handler.Auth.Detail))
			auth.POST("/update-password", authorize, core.HandlerFunc(handler.Auth.UpdatePassword))
			auth.POST("/update-detail", authorize, core.HandlerFunc(handler.Auth.UpdateDetail))
		}

		admins := v1.Group("/admin", authorize)
		{
			admins.POST("/list", core.HandlerFunc(handler.Admin.List))
			admins.POST("/create", core.HandlerFunc(handler.Admin.Create))
			admins.POST("/update-status", core.HandlerFunc(handler.Admin.UpdateStatus))
			admins.POST("/reset-password", core.HandlerFunc(handler.Admin.ResetPassword))
		}

		role := v1.Group("/role", authorize)
		{
			role.POST("/list", core.HandlerFunc(handler.Role.List))
			role.POST("/create", core.HandlerFunc(handler.Role.Create))
			role.POST("/update", core.HandlerFunc(handler.Role.Update))
		}

		resource := v1.Group("/resource", authorize)
		{
			resource.POST("/list", core.HandlerFunc(handler.Resource.List))
			resource.POST("/create", core.HandlerFunc(handler.Resource.Create))
			resource.POST("/update", core.HandlerFunc(handler.Resource.Update))
			resource.POST("/delete", core.HandlerFunc(handler.Resource.Delete))
		}

		menu := v1.Group("/menu", authorize)
		{
			menu.POST("/list", core.HandlerFunc(handler.Menu.List))
			menu.POST("/create", core.HandlerFunc(handler.Menu.Create))
			menu.POST("/update", core.HandlerFunc(handler.Menu.Update))
			menu.POST("/delete", core.HandlerFunc(handler.Menu.Delete))
			menu.POST("/user", core.HandlerFunc(handler.Menu.GetUserMenus))
		}

		totp := v1.Group("/totp", authorize)
		{
			totp.POST("/status", core.HandlerFunc(handler.Totp.Status))
			totp.POST("/init", core.HandlerFunc(handler.Totp.Init))
			totp.POST("/submit", core.HandlerFunc(handler.Totp.Submit))
			totp.POST("/qrcode", core.HandlerFunc(handler.Totp.Qrcode))
			totp.POST("/close", core.HandlerFunc(handler.Totp.Close))
		}
	}
}
