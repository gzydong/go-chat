package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/handler"
	"go-chat/app/http/middleware"
	"go-chat/config"
)

// RegisterApiRoute 注册 API 路由
func RegisterApiRoute(conf *config.Config, router *gin.Engine, handler *handler.Handler) {
	// 授权验证中间件
	authorize := middleware.JwtAuth(conf, "api")

	group := router.Group("/api/v1")
	{
		common := group.Group("/common")
		{
			common.POST("/sms-code", handler.Common.SmsCode)
		}

		// 授权相关分组
		auth := group.Group("/auth")
		{
			auth.POST("/login", handler.Auth.Login)
			auth.POST("/register", handler.Auth.Register)
			auth.POST("/refresh", authorize, handler.Auth.Refresh)
			auth.POST("/logout", authorize, handler.Auth.Logout)
			auth.POST("/forget", handler.Auth.Forget)
		}

		// 用户相关分组
		user := group.Group("/user").Use(authorize)
		{
			user.GET("/detail", handler.User.Detail)
			user.POST("/change/detail", handler.User.ChangeDetail)
			user.POST("/change/password", handler.User.ChangePassword)
			user.POST("/change/mobile", handler.User.ChangeMobile)
			user.POST("/change/email", handler.User.ChangeEmail)
		}

		talkMsg := group.Group("/talk/message").Use(authorize)
		{
			talkMsg.POST("/text", handler.TalkMessage.Text)
			talkMsg.POST("/code", handler.TalkMessage.Code)
			talkMsg.POST("/image", handler.TalkMessage.Image)
			talkMsg.POST("/file", handler.TalkMessage.File)
			talkMsg.POST("/emoticon", handler.TalkMessage.Emoticon)
			talkMsg.POST("/forward", handler.TalkMessage.Forward)
			talkMsg.POST("/card", handler.TalkMessage.Card)
			talkMsg.POST("/collect", handler.TalkMessage.Collect)
			talkMsg.POST("/revoke", handler.TalkMessage.Revoke)
			talkMsg.POST("/delete", handler.TalkMessage.Delete)
			talkMsg.POST("/vote", handler.TalkMessage.Vote)
			talkMsg.POST("/handleVote", handler.TalkMessage.HandleVote)
		}

		download := group.Group("/download").Use(authorize)
		{
			download.GET("/chat/file", handler.Download.ArticleAnnex)
		}

		emoticon := group.Group("/emoticon").Use(authorize)
		{
			emoticon.GET("/list", handler.Emoticon.CollectList)
			emoticon.POST("/upload", handler.Emoticon.Upload)
			emoticon.GET("/system", handler.Emoticon.SystemList)
			emoticon.POST("/set-user-emoticon", handler.Emoticon.SetSystemEmoticon)
			emoticon.POST("/del-collect-emoticon", handler.Emoticon.DeleteCollect)
		}

		upload := group.Group("/upload")
		{
			upload.POST("/index", handler.Upload.Index)
		}
	}
}
