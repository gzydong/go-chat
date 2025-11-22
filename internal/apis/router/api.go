package router

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	web2 "github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/apis/handler/web"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/jwtutil"
)

// RegisterWebRoute 注册 Web 路由
func RegisterWebRoute(secret string, router *gin.Engine, handler *web.Handler, storage middleware.IStorage) {
	// 授权验证中间件
	authorize := middleware.NewJwtMiddleware[entity.WebClaims](
		[]byte(secret), storage,
		func(ctx context.Context, claims *jwtutil.JwtClaims[entity.WebClaims]) error {
			if claims.RegisteredClaims.Issuer != entity.JwtIssuerWeb {
				return errors.New("授权异常，请登录后操作")
			}

			user, err := handler.UserRepo.FindById(ctx, claims.Metadata.UserId)
			if err != nil {
				return errors.New("授权异常，请登录后操作")
			}

			if user.IsDisabled() {
				return entity.ErrAccountDisabled
			}

			return nil
		},
	)

	api := router.Group("/").Use(authorize)

	resp := &Interceptor{}

	web2.RegisterAuthHandler(router, resp, handler.V1.Auth)
	web2.RegisterCommonHandler(router, resp, handler.V1.Common)
	web2.RegisterUserHandler(api, resp, handler.V1.User)
	web2.RegisterEmoticonHandler(api, resp, handler.V1.Emoticon)
	web2.RegisterOrganizeHandler(api, resp, handler.V1.Organize)
	web2.RegisterArticleClassHandler(api, resp, handler.V1.ArticleClass)
	web2.RegisterArticleHandler(api, resp, handler.V1.Article)
	web2.RegisterArticleAnnexHandler(api, resp, handler.V1.ArticleAnnex)
	web2.RegisterContactHandler(api, resp, handler.V1.Contact)
	web2.RegisterContactApplyHandler(api, resp, handler.V1.ContactApply)
	web2.RegisterContactGroupHandler(api, resp, handler.V1.ContactGroup)
	web2.RegisterTalkHandler(api, resp, handler.V1.Talk)
	web2.RegisterGroupHandler(api, resp, handler.V1.Group)
	web2.RegisterGroupApplyHandler(api, resp, handler.V1.GroupApply)
	web2.RegisterGroupVoteHandler(api, resp, handler.V1.GroupVote)
	web2.RegisterGroupNoticeHandler(api, resp, handler.V1.GroupNotice)
	web2.RegisterMessageHandler(api, resp, handler.V1.TalkMessage)

	registerCustomApiRouter(resp, api, handler)
}

func registerCustomApiRouter(resp *Interceptor, api gin.IRoutes, handler *web.Handler) {
	api.POST("/api/v1/emoticon/customize/upload", HandlerFunc(resp, func(c *gin.Context) (any, error) {
		return handler.V1.Emoticon.Upload(c, &web2.EmoticonUploadRequest{})
	}))

	api.POST("/api/v1/article-annex/upload", HandlerFunc(resp, func(c *gin.Context) (any, error) {
		return handler.V1.ArticleAnnex.Upload(c, &web2.ArticleAnnexUploadRequest{})
	}))

	api.GET("/api/v1/article-annex/download", func(c *gin.Context) {
		_, err := handler.V1.ArticleAnnex.Download(c, nil)
		if err != nil {
			resp.Error(c, err)
		}
	})

	api.POST("/api/v1/upload/media-file", HandlerFunc(resp, func(c *gin.Context) (any, error) {
		return handler.V1.Upload.Image(c)
	}))

	api.POST("/api/v1/upload/multipart", HandlerFunc(resp, func(c *gin.Context) (any, error) {
		return handler.V1.Upload.MultipartUpload(c)
	}))

	api.POST("/api/v1/upload/init-multipart", HandlerFunc(resp, func(c *gin.Context) (any, error) {
		return handler.V1.Upload.InitiateMultipart(c)
	}))

	api.GET("/api/v1/talk/file-download", func(c *gin.Context) {
		if err := handler.V1.TalkMessage.Download(c); err != nil {
			resp.Error(c, err)
		}
	})

	api.POST("/api/v1/message/send", HandlerFunc(resp, func(c *gin.Context) (any, error) {
		return handler.V1.Message.Send(c)
	}))
}
