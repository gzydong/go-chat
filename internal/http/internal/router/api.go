package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwtutil"
)

// RegisterApiRoute 注册 API 路由
func RegisterApiRoute(conf *config.Config, router *gin.Engine, handler *handler.ApiHandler, session *cache.Session) {
	// 授权验证中间件
	authorize := jwtutil.Auth(conf.Jwt.Secret, "api", session)

	// v1 接口
	v1 := router.Group("/api/v1")
	{
		test := v1.Group("/test")
		{
			test.GET("/success", ichat.NewHandlerFunc(handler.Test.Success))
			test.GET("/raw", ichat.NewHandlerFunc(handler.Test.Raw))
			test.GET("/error", ichat.NewHandlerFunc(handler.Test.Error))
			test.GET("/invalid", ichat.NewHandlerFunc(handler.Test.Invalid))
		}

		common := v1.Group("/common")
		{
			common.POST("/sms-code", ichat.NewHandlerFunc(handler.Common.SmsCode))
			common.POST("/email-code", authorize, ichat.NewHandlerFunc(handler.Common.EmailCode))
			common.GET("/setting", authorize, ichat.NewHandlerFunc(handler.Common.Setting))
		}

		// 授权相关分组
		auth := v1.Group("/auth")
		{
			auth.POST("/login", ichat.NewHandlerFunc(handler.Auth.Login))                // 登录
			auth.POST("/register", ichat.NewHandlerFunc(handler.Auth.Register))          // 注册
			auth.POST("/refresh", authorize, ichat.NewHandlerFunc(handler.Auth.Refresh)) // 刷新 Token
			auth.POST("/logout", authorize, ichat.NewHandlerFunc(handler.Auth.Logout))   // 退出登录
			auth.POST("/forget", ichat.NewHandlerFunc(handler.Auth.Forget))              // 找回密码
		}

		// 用户相关分组
		user := v1.Group("/users").Use(authorize)
		{
			user.GET("/detail", ichat.HandlerFunc(handler.User.Detail))                   // 获取个人信息
			user.GET("/setting", ichat.HandlerFunc(handler.User.Setting))                 // 获取个人信息
			user.POST("/change/detail", ichat.HandlerFunc(handler.User.ChangeDetail))     // 修改用户信息
			user.POST("/change/password", ichat.HandlerFunc(handler.User.ChangePassword)) // 修改用户密码
			user.POST("/change/mobile", ichat.HandlerFunc(handler.User.ChangeMobile))     // 修改用户手机号
			user.POST("/change/email", ichat.HandlerFunc(handler.User.ChangeEmail))       // 修改用户邮箱
		}

		contact := v1.Group("/contact").Use(authorize)
		{
			contact.GET("/list", ichat.HandlerFunc(handler.Contact.List))               // 联系人列表
			contact.GET("/search", ichat.HandlerFunc(handler.Contact.Search))           // 搜索联系人
			contact.GET("/detail", ichat.HandlerFunc(handler.Contact.Detail))           // 搜索联系人
			contact.POST("/delete", ichat.HandlerFunc(handler.Contact.Delete))          // 删除联系人
			contact.POST("/edit-remark", ichat.HandlerFunc(handler.Contact.EditRemark)) // 编辑联系人备注

			// 联系人申请相关
			contact.GET("/apply/records", ichat.HandlerFunc(handler.ContactsApply.List))              // 联系人申请列表
			contact.POST("/apply/create", ichat.HandlerFunc(handler.ContactsApply.Create))            // 添加联系人申请
			contact.POST("/apply/accept", ichat.HandlerFunc(handler.ContactsApply.Accept))            // 同意人申请列表
			contact.POST("/apply/decline", ichat.HandlerFunc(handler.ContactsApply.Decline))          // 拒绝人申请列表
			contact.GET("/apply/unread-num", ichat.HandlerFunc(handler.ContactsApply.ApplyUnreadNum)) // 联系人申请未读数
		}

		// 聊天群相关分组
		userGroup := v1.Group("/group").Use(authorize)
		{
			userGroup.GET("/list", ichat.HandlerFunc(handler.Group.GetGroups))            // 群组列表
			userGroup.GET("/overt/list", ichat.HandlerFunc(handler.Group.OvertList))      // 公开群组列表
			userGroup.GET("/detail", ichat.HandlerFunc(handler.Group.Detail))             // 群组详情
			userGroup.POST("/create", ichat.HandlerFunc(handler.Group.Create))            // 创建群组
			userGroup.POST("/dismiss", ichat.HandlerFunc(handler.Group.Dismiss))          // 解散群组
			userGroup.POST("/invite", ichat.HandlerFunc(handler.Group.Invite))            // 邀请加入群组
			userGroup.POST("/secede", ichat.HandlerFunc(handler.Group.SignOut))           // 退出群组
			userGroup.POST("/setting", ichat.HandlerFunc(handler.Group.Setting))          // 设置群组信息
			userGroup.POST("/handover", ichat.HandlerFunc(handler.Group.Handover))        // 群主转让
			userGroup.POST("/assign-admin", ichat.HandlerFunc(handler.Group.AssignAdmin)) // 分配管理员
			userGroup.POST("/no-speak", ichat.HandlerFunc(handler.Group.NoSpeak))         // 修改禁言状态

			// 群成员相关
			userGroup.GET("/member/list", ichat.HandlerFunc(handler.Group.GetMembers))          // 群成员列表
			userGroup.GET("/member/invites", ichat.HandlerFunc(handler.Group.GetInviteFriends)) // 群成员列表
			userGroup.POST("/member/remove", ichat.HandlerFunc(handler.Group.RemoveMembers))    // 移出指定群成员
			userGroup.POST("/member/remark", ichat.HandlerFunc(handler.Group.EditRemark))       // 设置群名片

			// 群公告相关
			userGroup.GET("/notice/list", ichat.HandlerFunc(handler.GroupNotice.List))             // 群公告列表
			userGroup.POST("/notice/edit", ichat.HandlerFunc(handler.GroupNotice.CreateAndUpdate)) // 添加或编辑群公告
			userGroup.POST("/notice/delete", ichat.HandlerFunc(handler.GroupNotice.Delete))        // 删除群公告

			// 群申请
			userGroup.POST("/apply/create", ichat.HandlerFunc(handler.GroupApply.Create)) // 提交入群申请
			userGroup.POST("/apply/delete", ichat.HandlerFunc(handler.GroupApply.Delete)) // 申请入群申请
			userGroup.POST("/apply/agree", ichat.HandlerFunc(handler.GroupApply.Agree))   // 同意入群申请
			userGroup.GET("/apply/list", ichat.HandlerFunc(handler.GroupApply.List))      // 入群申请列表
		}

		talk := v1.Group("/talk").Use(authorize)
		{
			talk.GET("/list", ichat.HandlerFunc(handler.Talk.List))                                   // 会话列表
			talk.POST("/create", ichat.HandlerFunc(handler.Talk.Create))                              // 创建会话
			talk.POST("/delete", ichat.HandlerFunc(handler.Talk.Delete))                              // 删除会话
			talk.POST("/topping", ichat.HandlerFunc(handler.Talk.Top))                                // 置顶会话
			talk.POST("/disturb", ichat.HandlerFunc(handler.Talk.Disturb))                            // 会话免打扰
			talk.GET("/records", ichat.HandlerFunc(handler.TalkRecords.GetRecords))                   // 会话面板记录
			talk.GET("/records/history", ichat.HandlerFunc(handler.TalkRecords.SearchHistoryRecords)) // 历史会话记录
			talk.GET("/records/forward", ichat.HandlerFunc(handler.TalkRecords.GetForwardRecords))    // 会话转发记录
			talk.GET("/records/file/download", ichat.HandlerFunc(handler.TalkRecords.Download))       // 会话转发记录
			talk.POST("/unread/clear", ichat.HandlerFunc(handler.Talk.ClearUnreadMessage))            // 清除会话未读数
		}

		talkMsg := v1.Group("/talk/message").Use(authorize)
		{
			talkMsg.POST("/text", ichat.HandlerFunc(handler.TalkMessage.Text))              // 发送文本消息
			talkMsg.POST("/code", ichat.HandlerFunc(handler.TalkMessage.Code))              // 发送代码消息
			talkMsg.POST("/image", ichat.HandlerFunc(handler.TalkMessage.Image))            // 发送图片消息
			talkMsg.POST("/file", ichat.HandlerFunc(handler.TalkMessage.File))              // 发送文件消息
			talkMsg.POST("/emoticon", ichat.HandlerFunc(handler.TalkMessage.Emoticon))      // 发送表情包消息
			talkMsg.POST("/forward", ichat.HandlerFunc(handler.TalkMessage.Forward))        // 发送转发消息
			talkMsg.POST("/card", ichat.HandlerFunc(handler.TalkMessage.Card))              // 发送用户名片
			talkMsg.POST("/location", ichat.HandlerFunc(handler.TalkMessage.Location))      // 发送位置消息
			talkMsg.POST("/collect", ichat.HandlerFunc(handler.TalkMessage.Collect))        // 收藏会话表情图片
			talkMsg.POST("/revoke", ichat.HandlerFunc(handler.TalkMessage.Revoke))          // 撤销聊天消息
			talkMsg.POST("/delete", ichat.HandlerFunc(handler.TalkMessage.Delete))          // 删除聊天消息
			talkMsg.POST("/vote", ichat.HandlerFunc(handler.TalkMessage.Vote))              // 发送投票消息
			talkMsg.POST("/vote/handle", ichat.HandlerFunc(handler.TalkMessage.HandleVote)) // 投票消息处理
		}

		emoticon := v1.Group("/emoticon").Use(authorize)
		{
			emoticon.GET("/list", ichat.HandlerFunc(handler.Emoticon.CollectList))                // 表情包列表
			emoticon.POST("/customize/create", ichat.HandlerFunc(handler.Emoticon.Upload))        // 添加自定义表情
			emoticon.POST("/customize/delete", ichat.HandlerFunc(handler.Emoticon.DeleteCollect)) // 删除自定义表情

			// 系統表情包
			emoticon.GET("/system/list", ichat.HandlerFunc(handler.Emoticon.SystemList))            // 系统表情包列表
			emoticon.POST("/system/install", ichat.HandlerFunc(handler.Emoticon.SetSystemEmoticon)) // 添加或移除系统表情包
		}

		upload := v1.Group("/upload").Use(authorize)
		{
			upload.POST("/avatar", ichat.HandlerFunc(handler.Upload.Avatar))
			upload.POST("/multipart/initiate", ichat.HandlerFunc(handler.Upload.InitiateMultipart))
			upload.POST("/multipart", ichat.HandlerFunc(handler.Upload.MultipartUpload))
		}

		note := v1.Group("/note").Use(authorize)
		{
			// 文章相关
			note.GET("/article/list", ichat.HandlerFunc(handler.Article.List))
			note.POST("/article/editor", ichat.HandlerFunc(handler.Article.Edit))
			note.GET("/article/detail", ichat.HandlerFunc(handler.Article.Detail))
			note.POST("/article/delete", ichat.HandlerFunc(handler.Article.Delete))
			note.POST("/article/upload/image", ichat.HandlerFunc(handler.Article.Upload))
			note.POST("/article/recover", ichat.HandlerFunc(handler.Article.Recover))
			note.POST("/article/move", ichat.HandlerFunc(handler.Article.Move))
			note.POST("/article/asterisk", ichat.HandlerFunc(handler.Article.Asterisk))
			note.POST("/article/tag", ichat.HandlerFunc(handler.Article.Tag))
			note.POST("/article/forever/delete", ichat.HandlerFunc(handler.Article.ForeverDelete))

			// 文章分类
			note.GET("/class/list", ichat.HandlerFunc(handler.ArticleClass.List))
			note.POST("/class/editor", ichat.HandlerFunc(handler.ArticleClass.Edit))
			note.POST("/class/delete", ichat.HandlerFunc(handler.ArticleClass.Delete))
			note.POST("/class/sort", ichat.HandlerFunc(handler.ArticleClass.Sort))

			// 文章标签
			note.GET("/tag/list", ichat.HandlerFunc(handler.ArticleTag.List))
			note.POST("/tag/editor", ichat.HandlerFunc(handler.ArticleTag.Edit))
			note.POST("/tag/delete", ichat.HandlerFunc(handler.ArticleTag.Delete))

			// 文章附件
			note.POST("/annex/upload", ichat.HandlerFunc(handler.ArticleAnnex.Upload))
			note.POST("/annex/delete", ichat.HandlerFunc(handler.ArticleAnnex.Delete))
			note.POST("/annex/recover", ichat.HandlerFunc(handler.ArticleAnnex.Recover))
			note.POST("/annex/forever/delete", ichat.HandlerFunc(handler.ArticleAnnex.ForeverDelete))
			note.GET("/annex/recover/list", ichat.HandlerFunc(handler.ArticleAnnex.RecoverList))
			note.GET("/annex/download", ichat.HandlerFunc(handler.ArticleAnnex.Download))
		}

		organize := v1.Group("/organize").Use(authorize)
		{
			organize.GET("/department/all", ichat.NewHandlerFunc(handler.Organize.DepartmentList))
			organize.GET("/personnel/all", ichat.NewHandlerFunc(handler.Organize.PersonnelList))
		}
	}

	// v2 接口
	v2 := router.Group("/api/v2")
	{
		v2.GET("/test", func(context *gin.Context) {
			context.JSON(200, entity.H{"message": "success"})
		})
	}
}
