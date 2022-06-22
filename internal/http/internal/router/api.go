package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/pkg/ginutil"
	"go-chat/internal/pkg/jwtutil"
)

// RegisterApiRoute 注册 API 路由
func RegisterApiRoute(conf *config.Config, router *gin.Engine, handler *handler.ApiHandler, session *cache.Session) {
	// 授权验证中间件
	authorize := jwtutil.Auth(conf.Jwt.Secret, "api", session)

	// v1 接口
	v1 := router.Group("/api/v1")
	{
		common := v1.Group("/common")
		{
			common.POST("/sms-code", ginutil.HandlerFunc(handler.Common.SmsCode))
			common.POST("/email-code", authorize, ginutil.HandlerFunc(handler.Common.EmailCode))
			common.GET("/setting", authorize, ginutil.HandlerFunc(handler.Common.Setting))
		}

		// 授权相关分组
		auth := v1.Group("/auth")
		{
			auth.POST("/login", ginutil.HandlerFunc(handler.Auth.Login))                // 登录
			auth.POST("/register", ginutil.HandlerFunc(handler.Auth.Register))          // 注册
			auth.POST("/refresh", authorize, ginutil.HandlerFunc(handler.Auth.Refresh)) // 刷新 Token
			auth.POST("/logout", authorize, ginutil.HandlerFunc(handler.Auth.Logout))   // 退出登录
			auth.POST("/forget", ginutil.HandlerFunc(handler.Auth.Forget))              // 找回密码
		}

		// 用户相关分组
		user := v1.Group("/users").Use(authorize)
		{
			user.GET("/detail", ginutil.HandlerFunc(handler.User.Detail))                   // 获取个人信息
			user.GET("/setting", ginutil.HandlerFunc(handler.User.Setting))                 // 获取个人信息
			user.POST("/change/detail", ginutil.HandlerFunc(handler.User.ChangeDetail))     // 修改用户信息
			user.POST("/change/password", ginutil.HandlerFunc(handler.User.ChangePassword)) // 修改用户密码
			user.POST("/change/mobile", ginutil.HandlerFunc(handler.User.ChangeMobile))     // 修改用户手机号
			user.POST("/change/email", ginutil.HandlerFunc(handler.User.ChangeEmail))       // 修改用户邮箱
		}

		contact := v1.Group("/contact").Use(authorize)
		{
			contact.GET("/list", ginutil.HandlerFunc(handler.Contact.List))               // 联系人列表
			contact.GET("/search", ginutil.HandlerFunc(handler.Contact.Search))           // 搜索联系人
			contact.GET("/detail", ginutil.HandlerFunc(handler.Contact.Detail))           // 搜索联系人
			contact.POST("/delete", ginutil.HandlerFunc(handler.Contact.Delete))          // 删除联系人
			contact.POST("/edit-remark", ginutil.HandlerFunc(handler.Contact.EditRemark)) // 编辑联系人备注

			// 联系人申请相关
			contact.GET("/apply/records", ginutil.HandlerFunc(handler.ContactsApply.List))              // 联系人申请列表
			contact.POST("/apply/create", ginutil.HandlerFunc(handler.ContactsApply.Create))            // 添加联系人申请
			contact.POST("/apply/accept", ginutil.HandlerFunc(handler.ContactsApply.Accept))            // 同意人申请列表
			contact.POST("/apply/decline", ginutil.HandlerFunc(handler.ContactsApply.Decline))          // 拒绝人申请列表
			contact.GET("/apply/unread-num", ginutil.HandlerFunc(handler.ContactsApply.ApplyUnreadNum)) // 联系人申请未读数
		}

		// 聊天群相关分组
		userGroup := v1.Group("/group").Use(authorize)
		{
			userGroup.GET("/list", ginutil.HandlerFunc(handler.Group.GetGroups))            // 群组列表
			userGroup.GET("/overt/list", ginutil.HandlerFunc(handler.Group.OvertList))      // 公开群组列表
			userGroup.GET("/detail", ginutil.HandlerFunc(handler.Group.Detail))             // 群组详情
			userGroup.POST("/create", ginutil.HandlerFunc(handler.Group.Create))            // 创建群组
			userGroup.POST("/dismiss", ginutil.HandlerFunc(handler.Group.Dismiss))          // 解散群组
			userGroup.POST("/invite", ginutil.HandlerFunc(handler.Group.Invite))            // 邀请加入群组
			userGroup.POST("/secede", ginutil.HandlerFunc(handler.Group.SignOut))           // 退出群组
			userGroup.POST("/setting", ginutil.HandlerFunc(handler.Group.Setting))          // 设置群组信息
			userGroup.POST("/handover", ginutil.HandlerFunc(handler.Group.Handover))        // 群主转让
			userGroup.POST("/assign-admin", ginutil.HandlerFunc(handler.Group.AssignAdmin)) // 分配管理员
			userGroup.POST("/no-speak", ginutil.HandlerFunc(handler.Group.NoSpeak))         // 修改禁言状态

			// 群成员相关
			userGroup.GET("/member/list", ginutil.HandlerFunc(handler.Group.GetMembers))          // 群成员列表
			userGroup.GET("/member/invites", ginutil.HandlerFunc(handler.Group.GetInviteFriends)) // 群成员列表
			userGroup.POST("/member/remove", ginutil.HandlerFunc(handler.Group.RemoveMembers))    // 移出指定群成员
			userGroup.POST("/member/remark", ginutil.HandlerFunc(handler.Group.EditRemark))       // 设置群名片

			// 群公告相关
			userGroup.GET("/notice/list", ginutil.HandlerFunc(handler.GroupNotice.List))             // 群公告列表
			userGroup.POST("/notice/edit", ginutil.HandlerFunc(handler.GroupNotice.CreateAndUpdate)) // 添加或编辑群公告
			userGroup.POST("/notice/delete", ginutil.HandlerFunc(handler.GroupNotice.Delete))        // 删除群公告

			// 群申请
			userGroup.POST("/apply/create", ginutil.HandlerFunc(handler.GroupApply.Create)) // 提交入群申请
			userGroup.POST("/apply/delete", ginutil.HandlerFunc(handler.GroupApply.Delete)) // 申请入群申请
			userGroup.POST("/apply/agree", ginutil.HandlerFunc(handler.GroupApply.Agree))   // 同意入群申请
			userGroup.GET("/apply/list", ginutil.HandlerFunc(handler.GroupApply.List))      // 入群申请列表
		}

		talk := v1.Group("/talk").Use(authorize)
		{
			talk.GET("/list", ginutil.HandlerFunc(handler.Talk.List))                                   // 会话列表
			talk.POST("/create", ginutil.HandlerFunc(handler.Talk.Create))                              // 创建会话
			talk.POST("/delete", ginutil.HandlerFunc(handler.Talk.Delete))                              // 删除会话
			talk.POST("/topping", ginutil.HandlerFunc(handler.Talk.Top))                                // 置顶会话
			talk.POST("/disturb", ginutil.HandlerFunc(handler.Talk.Disturb))                            // 会话免打扰
			talk.GET("/records", ginutil.HandlerFunc(handler.TalkRecords.GetRecords))                   // 会话面板记录
			talk.GET("/records/history", ginutil.HandlerFunc(handler.TalkRecords.SearchHistoryRecords)) // 历史会话记录
			talk.GET("/records/forward", ginutil.HandlerFunc(handler.TalkRecords.GetForwardRecords))    // 会话转发记录
			talk.GET("/records/file/download", ginutil.HandlerFunc(handler.TalkRecords.Download))       // 会话转发记录
			talk.POST("/unread/clear", ginutil.HandlerFunc(handler.Talk.ClearUnreadMessage))            // 清除会话未读数
		}

		talkMsg := v1.Group("/talk/message").Use(authorize)
		{
			talkMsg.POST("/text", ginutil.HandlerFunc(handler.TalkMessage.Text))              // 发送文本消息
			talkMsg.POST("/code", ginutil.HandlerFunc(handler.TalkMessage.Code))              // 发送代码消息
			talkMsg.POST("/image", ginutil.HandlerFunc(handler.TalkMessage.Image))            // 发送图片消息
			talkMsg.POST("/file", ginutil.HandlerFunc(handler.TalkMessage.File))              // 发送文件消息
			talkMsg.POST("/emoticon", ginutil.HandlerFunc(handler.TalkMessage.Emoticon))      // 发送表情包消息
			talkMsg.POST("/forward", ginutil.HandlerFunc(handler.TalkMessage.Forward))        // 发送转发消息
			talkMsg.POST("/card", ginutil.HandlerFunc(handler.TalkMessage.Card))              // 发送用户名片
			talkMsg.POST("/location", ginutil.HandlerFunc(handler.TalkMessage.Location))      // 发送位置消息
			talkMsg.POST("/collect", ginutil.HandlerFunc(handler.TalkMessage.Collect))        // 收藏会话表情图片
			talkMsg.POST("/revoke", ginutil.HandlerFunc(handler.TalkMessage.Revoke))          // 撤销聊天消息
			talkMsg.POST("/delete", ginutil.HandlerFunc(handler.TalkMessage.Delete))          // 删除聊天消息
			talkMsg.POST("/vote", ginutil.HandlerFunc(handler.TalkMessage.Vote))              // 发送投票消息
			talkMsg.POST("/vote/handle", ginutil.HandlerFunc(handler.TalkMessage.HandleVote)) // 投票消息处理
		}

		emoticon := v1.Group("/emoticon").Use(authorize)
		{
			emoticon.GET("/list", ginutil.HandlerFunc(handler.Emoticon.CollectList))                // 表情包列表
			emoticon.POST("/customize/create", ginutil.HandlerFunc(handler.Emoticon.Upload))        // 添加自定义表情
			emoticon.POST("/customize/delete", ginutil.HandlerFunc(handler.Emoticon.DeleteCollect)) // 删除自定义表情

			// 系統表情包
			emoticon.GET("/system/list", ginutil.HandlerFunc(handler.Emoticon.SystemList))            // 系统表情包列表
			emoticon.POST("/system/install", ginutil.HandlerFunc(handler.Emoticon.SetSystemEmoticon)) // 添加或移除系统表情包
		}

		upload := v1.Group("/upload").Use(authorize)
		{
			upload.POST("/avatar", ginutil.HandlerFunc(handler.Upload.Avatar))
			upload.POST("/multipart/initiate", ginutil.HandlerFunc(handler.Upload.InitiateMultipart))
			upload.POST("/multipart", ginutil.HandlerFunc(handler.Upload.MultipartUpload))
		}

		note := v1.Group("/note").Use(authorize)
		{
			// 文章相关
			note.GET("/article/list", ginutil.HandlerFunc(handler.Article.List))
			note.POST("/article/editor", ginutil.HandlerFunc(handler.Article.Edit))
			note.GET("/article/detail", ginutil.HandlerFunc(handler.Article.Detail))
			note.POST("/article/delete", ginutil.HandlerFunc(handler.Article.Delete))
			note.POST("/article/upload/image", ginutil.HandlerFunc(handler.Article.Upload))
			note.POST("/article/recover", ginutil.HandlerFunc(handler.Article.Recover))
			note.POST("/article/move", ginutil.HandlerFunc(handler.Article.Move))
			note.POST("/article/asterisk", ginutil.HandlerFunc(handler.Article.Asterisk))
			note.POST("/article/tag", ginutil.HandlerFunc(handler.Article.Tag))
			note.POST("/article/forever/delete", ginutil.HandlerFunc(handler.Article.ForeverDelete))

			// 文章分类
			note.GET("/class/list", ginutil.HandlerFunc(handler.ArticleClass.List))
			note.POST("/class/editor", ginutil.HandlerFunc(handler.ArticleClass.Edit))
			note.POST("/class/delete", ginutil.HandlerFunc(handler.ArticleClass.Delete))
			note.POST("/class/sort", ginutil.HandlerFunc(handler.ArticleClass.Sort))

			// 文章标签
			note.GET("/tag/list", ginutil.HandlerFunc(handler.ArticleTag.List))
			note.POST("/tag/editor", ginutil.HandlerFunc(handler.ArticleTag.Edit))
			note.POST("/tag/delete", ginutil.HandlerFunc(handler.ArticleTag.Delete))

			// 文章附件
			note.POST("/annex/upload", ginutil.HandlerFunc(handler.ArticleAnnex.Upload))
			note.POST("/annex/delete", ginutil.HandlerFunc(handler.ArticleAnnex.Delete))
			note.POST("/annex/recover", ginutil.HandlerFunc(handler.ArticleAnnex.Recover))
			note.POST("/annex/forever/delete", ginutil.HandlerFunc(handler.ArticleAnnex.ForeverDelete))
			note.GET("/annex/recover/list", ginutil.HandlerFunc(handler.ArticleAnnex.RecoverList))
			note.GET("/annex/download", ginutil.HandlerFunc(handler.ArticleAnnex.Download))
		}

		organize := v1.Group("/organize").Use(authorize)
		{
			organize.GET("/department/all", ginutil.HandlerFunc(handler.Organize.DepartmentList))
			organize.GET("/personnel/all", ginutil.HandlerFunc(handler.Organize.PersonnelList))
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
