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
			common.POST("/email-code", authorize, handler.Common.EmailCode)
			common.GET("/setting", authorize, handler.Common.Setting)
		}

		// 授权相关分组
		auth := group.Group("/auth")
		{
			auth.POST("/login", handler.Auth.Login)                // 登录接口
			auth.POST("/register", handler.Auth.Register)          // 注册接口
			auth.POST("/refresh", authorize, handler.Auth.Refresh) // Token 刷新接口
			auth.POST("/logout", authorize, handler.Auth.Logout)   // 退出登录接口
			auth.POST("/forget", handler.Auth.Forget)              // 找回密码接口
		}

		// 用户相关分组
		user := group.Group("/user").Use(authorize)
		{
			user.GET("/detail", handler.User.Detail)                   // 获取个人信息
			user.POST("/change/detail", handler.User.ChangeDetail)     // 修改个人信息接口
			user.POST("/change/password", handler.User.ChangePassword) // 修改密码接口
			user.POST("/change/mobile", handler.User.ChangeMobile)     // 修改手机号接口
			user.POST("/change/email", handler.User.ChangeEmail)       // 修改邮箱接口
		}

		// 聊天群相关分组
		userGroup := group.Group("/group").Use(authorize)
		{
			userGroup.POST("/create", handler.Group.Create)
			userGroup.POST("/dismiss", handler.Group.Dismiss)
			userGroup.POST("/invite", handler.Group.Invite)
			userGroup.POST("/secede", handler.Group.SignOut)
			userGroup.POST("/setting", handler.Group.Setting)
			userGroup.GET("/detail", handler.Group.Detail)
			userGroup.GET("/list", handler.Group.GetGroups)

			// 群成员相关
			userGroup.GET("/members", handler.Group.GetGroupMembers)          // 群成员列表
			userGroup.POST("/members/remove", handler.Group.RemoveMembers)    // 移出指定群成员
			userGroup.POST("/members/card", handler.Group.EditGroupRemarks)   // 设置用户群名片
			userGroup.POST("/invite-friends", handler.Group.GetInviteFriends) // 获取可邀请加入群组的好友列表

			// 群公告相关
			userGroup.GET("/notice/list", handler.GroupNotice.List)
			userGroup.POST("/notice/edit", handler.GroupNotice.CreateAndUpdate)
			userGroup.POST("/notice/delete", handler.GroupNotice.Delete)
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
			talkMsg.POST("/location", handler.TalkMessage.Location)
			talkMsg.POST("/collect", handler.TalkMessage.Collect)
			talkMsg.POST("/revoke", handler.TalkMessage.Revoke)
			talkMsg.POST("/delete", handler.TalkMessage.Delete)
			talkMsg.POST("/vote", handler.TalkMessage.Vote)
			talkMsg.POST("/vote/handle", handler.TalkMessage.HandleVote)
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
