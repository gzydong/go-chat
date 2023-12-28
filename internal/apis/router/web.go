package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/apis/handler/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/ichat/middleware"
)

// RegisterWebRoute 注册 Web 路由
func RegisterWebRoute(secret string, router *gin.Engine, handler *web.Handler, storage middleware.IStorage) {

	// 授权验证中间件
	authorize := middleware.Auth(secret, "api", storage)

	// v1 接口
	v1 := router.Group("/api/v1")
	{
		common := v1.Group("/common")
		{
			common.POST("/sms-code", ichat.HandlerFunc(handler.V1.Common.SmsCode))
			common.POST("/email-code", authorize, ichat.HandlerFunc(handler.V1.Common.EmailCode))
			common.GET("/setting", authorize, ichat.HandlerFunc(handler.V1.Common.Setting))
		}

		// 授权相关分组
		auth := v1.Group("/auth")
		{
			auth.POST("/login", ichat.HandlerFunc(handler.V1.Auth.Login))                // 登录
			auth.POST("/register", ichat.HandlerFunc(handler.V1.Auth.Register))          // 注册
			auth.POST("/refresh", authorize, ichat.HandlerFunc(handler.V1.Auth.Refresh)) // 刷新 Token
			auth.POST("/logout", authorize, ichat.HandlerFunc(handler.V1.Auth.Logout))   // 退出登录
			auth.POST("/forget", ichat.HandlerFunc(handler.V1.Auth.Forget))              // 找回密码
		}

		// 用户相关分组
		user := v1.Group("/users").Use(authorize)
		{
			user.GET("/detail", ichat.HandlerFunc(handler.V1.User.Detail))                   // 获取个人信息
			user.GET("/setting", ichat.HandlerFunc(handler.V1.User.Setting))                 // 获取个人信息
			user.POST("/change/detail", ichat.HandlerFunc(handler.V1.User.ChangeDetail))     // 修改用户信息
			user.POST("/change/password", ichat.HandlerFunc(handler.V1.User.ChangePassword)) // 修改用户密码
			user.POST("/change/mobile", ichat.HandlerFunc(handler.V1.User.ChangeMobile))     // 修改用户手机号
			user.POST("/change/email", ichat.HandlerFunc(handler.V1.User.ChangeEmail))       // 修改用户邮箱
		}

		contact := v1.Group("/contact").Use(authorize)
		{
			contact.GET("/list", ichat.HandlerFunc(handler.V1.Contact.List))             // 联系人列表
			contact.GET("/search", ichat.HandlerFunc(handler.V1.Contact.Search))         // 搜索联系人
			contact.GET("/detail", ichat.HandlerFunc(handler.V1.Contact.Detail))         // 搜索联系人
			contact.POST("/delete", ichat.HandlerFunc(handler.V1.Contact.Delete))        // 删除联系人
			contact.POST("/edit-remark", ichat.HandlerFunc(handler.V1.Contact.Remark))   // 编辑联系人备注
			contact.POST("/move-group", ichat.HandlerFunc(handler.V1.Contact.MoveGroup)) // 编辑联系人备注

			// 联系人申请相关
			contact.GET("/apply/records", ichat.HandlerFunc(handler.V1.ContactApply.List))              // 联系人申请列表
			contact.POST("/apply/create", ichat.HandlerFunc(handler.V1.ContactApply.Create))            // 添加联系人申请
			contact.POST("/apply/accept", ichat.HandlerFunc(handler.V1.ContactApply.Accept))            // 同意人申请列表
			contact.POST("/apply/decline", ichat.HandlerFunc(handler.V1.ContactApply.Decline))          // 拒绝人申请列表
			contact.GET("/apply/unread-num", ichat.HandlerFunc(handler.V1.ContactApply.ApplyUnreadNum)) // 联系人申请未读数

			// 联系人分组
			contact.GET("/group/list", ichat.HandlerFunc(handler.V1.ContactGroup.List))  // 联系人分组列表
			contact.POST("/group/save", ichat.HandlerFunc(handler.V1.ContactGroup.Save)) // 联系人分组排序
		}

		// 聊天群相关分组
		userGroup := v1.Group("/group").Use(authorize)
		{
			userGroup.GET("/list", ichat.HandlerFunc(handler.V1.Group.GroupList))            // 群组列表
			userGroup.GET("/overt/list", ichat.HandlerFunc(handler.V1.Group.OvertList))      // 公开群组列表
			userGroup.GET("/detail", ichat.HandlerFunc(handler.V1.Group.Detail))             // 群组详情
			userGroup.POST("/create", ichat.HandlerFunc(handler.V1.Group.Create))            // 创建群组
			userGroup.POST("/dismiss", ichat.HandlerFunc(handler.V1.Group.Dismiss))          // 解散群组
			userGroup.POST("/invite", ichat.HandlerFunc(handler.V1.Group.Invite))            // 邀请加入群组
			userGroup.POST("/secede", ichat.HandlerFunc(handler.V1.Group.SignOut))           // 退出群组
			userGroup.POST("/setting", ichat.HandlerFunc(handler.V1.Group.Setting))          // 设置群组信息
			userGroup.POST("/handover", ichat.HandlerFunc(handler.V1.Group.Handover))        // 群主转让
			userGroup.POST("/assign-admin", ichat.HandlerFunc(handler.V1.Group.AssignAdmin)) // 分配管理员
			userGroup.POST("/no-speak", ichat.HandlerFunc(handler.V1.Group.NoSpeak))         // 修改禁言状态
			userGroup.POST("/mute", ichat.HandlerFunc(handler.V1.Group.Mute))                // 修改禁言状态
			userGroup.POST("/overt", ichat.HandlerFunc(handler.V1.Group.Overt))              // 修改禁言状态

			// 群成员相关
			userGroup.GET("/member/list", ichat.HandlerFunc(handler.V1.Group.Members))               // 群成员列表
			userGroup.GET("/member/invites", ichat.HandlerFunc(handler.V1.Group.GetInviteFriends))   // 群成员列表
			userGroup.POST("/member/remove", ichat.HandlerFunc(handler.V1.Group.RemoveMembers))      // 移出指定群成员
			userGroup.POST("/member/remark", ichat.HandlerFunc(handler.V1.Group.UpdateMemberRemark)) // 设置群名片

			// 群公告相关
			userGroup.GET("/notice/list", ichat.HandlerFunc(handler.V1.GroupNotice.List))             // 群公告列表
			userGroup.POST("/notice/edit", ichat.HandlerFunc(handler.V1.GroupNotice.CreateAndUpdate)) // 添加或编辑群公告
			userGroup.POST("/notice/delete", ichat.HandlerFunc(handler.V1.GroupNotice.Delete))        // 删除群公告

			// 群申请
			userGroup.POST("/apply/create", ichat.HandlerFunc(handler.V1.GroupApply.Create))        // 提交入群申请
			userGroup.POST("/apply/agree", ichat.HandlerFunc(handler.V1.GroupApply.Agree))          // 同意入群申请
			userGroup.POST("/apply/decline", ichat.HandlerFunc(handler.V1.GroupApply.Decline))      // 拒绝入群申请
			userGroup.GET("/apply/list", ichat.HandlerFunc(handler.V1.GroupApply.List))             // 入群申请列表
			userGroup.GET("/apply/all", ichat.HandlerFunc(handler.V1.GroupApply.All))               // 入群申请列表
			userGroup.GET("/apply/unread", ichat.HandlerFunc(handler.V1.GroupApply.ApplyUnreadNum)) // 入群申请未读
		}

		talk := v1.Group("/talk").Use(authorize)
		{
			talk.GET("/list", ichat.HandlerFunc(handler.V1.Talk.List))                                   // 会话列表
			talk.POST("/create", ichat.HandlerFunc(handler.V1.Talk.Create))                              // 创建会话
			talk.POST("/delete", ichat.HandlerFunc(handler.V1.Talk.Delete))                              // 删除会话
			talk.POST("/topping", ichat.HandlerFunc(handler.V1.Talk.Top))                                // 置顶会话
			talk.POST("/disturb", ichat.HandlerFunc(handler.V1.Talk.Disturb))                            // 会话免打扰
			talk.GET("/records", ichat.HandlerFunc(handler.V1.TalkRecords.GetRecords))                   // 会话面板记录
			talk.GET("/records/history", ichat.HandlerFunc(handler.V1.TalkRecords.SearchHistoryRecords)) // 历史会话记录
			talk.GET("/records/forward", ichat.HandlerFunc(handler.V1.TalkRecords.GetForwardRecords))    // 会话转发记录
			talk.GET("/records/file/download", ichat.HandlerFunc(handler.V1.TalkRecords.Download))       // 会话转发记录
			talk.POST("/unread/clear", ichat.HandlerFunc(handler.V1.Talk.ClearUnreadMessage))            // 清除会话未读数
		}

		talkMsg := v1.Group("/talk/message").Use(authorize)
		{
			talkMsg.POST("/publish", ichat.HandlerFunc(handler.V1.Message.Publish))            // 发送文本消息
			talkMsg.POST("/file", ichat.HandlerFunc(handler.V1.TalkMessage.File))              // 发送文件消息
			talkMsg.POST("/collect", ichat.HandlerFunc(handler.V1.TalkMessage.Collect))        // 收藏会话表情图片
			talkMsg.POST("/revoke", ichat.HandlerFunc(handler.V1.TalkMessage.Revoke))          // 撤销聊天消息
			talkMsg.POST("/delete", ichat.HandlerFunc(handler.V1.TalkMessage.Delete))          // 删除聊天消息
			talkMsg.POST("/vote", ichat.HandlerFunc(handler.V1.TalkMessage.Vote))              // 发送投票消息
			talkMsg.POST("/vote/handle", ichat.HandlerFunc(handler.V1.TalkMessage.HandleVote)) // 投票消息处理
		}

		emoticon := v1.Group("/emoticon").Use(authorize)
		{
			emoticon.GET("/list", ichat.HandlerFunc(handler.V1.Emoticon.CollectList))                // 表情包列表
			emoticon.POST("/customize/create", ichat.HandlerFunc(handler.V1.Emoticon.Upload))        // 添加自定义表情
			emoticon.POST("/customize/delete", ichat.HandlerFunc(handler.V1.Emoticon.DeleteCollect)) // 删除自定义表情

			// 系統表情包
			emoticon.GET("/system/list", ichat.HandlerFunc(handler.V1.Emoticon.SystemList))            // 系统表情包列表
			emoticon.POST("/system/install", ichat.HandlerFunc(handler.V1.Emoticon.SetSystemEmoticon)) // 添加或移除系统表情包
		}

		upload := v1.Group("/upload").Use(authorize)
		{
			upload.POST("/avatar", ichat.HandlerFunc(handler.V1.Upload.Avatar))
			upload.POST("/image", ichat.HandlerFunc(handler.V1.Upload.Image))
			upload.POST("/multipart/initiate", ichat.HandlerFunc(handler.V1.Upload.InitiateMultipart))
			upload.POST("/multipart", ichat.HandlerFunc(handler.V1.Upload.MultipartUpload))
		}

		note := v1.Group("/note").Use(authorize)
		{
			// 文章相关
			note.GET("/article/list", ichat.HandlerFunc(handler.V1.Article.List))
			note.POST("/article/editor", ichat.HandlerFunc(handler.V1.Article.Edit))
			note.GET("/article/detail", ichat.HandlerFunc(handler.V1.Article.Detail))
			note.POST("/article/delete", ichat.HandlerFunc(handler.V1.Article.Delete))
			note.POST("/article/upload/image", ichat.HandlerFunc(handler.V1.Article.Upload))
			note.POST("/article/recover", ichat.HandlerFunc(handler.V1.Article.Recover))
			note.POST("/article/move", ichat.HandlerFunc(handler.V1.Article.Move))
			note.POST("/article/asterisk", ichat.HandlerFunc(handler.V1.Article.Asterisk))
			note.POST("/article/tag", ichat.HandlerFunc(handler.V1.Article.Tag))
			note.POST("/article/forever/delete", ichat.HandlerFunc(handler.V1.Article.ForeverDelete))

			// 文章分类
			note.GET("/class/list", ichat.HandlerFunc(handler.V1.ArticleClass.List))
			note.POST("/class/editor", ichat.HandlerFunc(handler.V1.ArticleClass.Edit))
			note.POST("/class/delete", ichat.HandlerFunc(handler.V1.ArticleClass.Delete))
			note.POST("/class/sort", ichat.HandlerFunc(handler.V1.ArticleClass.Sort))

			// 文章标签
			note.GET("/tag/list", ichat.HandlerFunc(handler.V1.ArticleTag.List))
			note.POST("/tag/editor", ichat.HandlerFunc(handler.V1.ArticleTag.Edit))
			note.POST("/tag/delete", ichat.HandlerFunc(handler.V1.ArticleTag.Delete))

			// 文章附件
			note.POST("/annex/upload", ichat.HandlerFunc(handler.V1.ArticleAnnex.Upload))
			note.POST("/annex/delete", ichat.HandlerFunc(handler.V1.ArticleAnnex.Delete))
			note.POST("/annex/recover", ichat.HandlerFunc(handler.V1.ArticleAnnex.Recover))
			note.POST("/annex/forever/delete", ichat.HandlerFunc(handler.V1.ArticleAnnex.ForeverDelete))
			note.GET("/annex/recover/list", ichat.HandlerFunc(handler.V1.ArticleAnnex.RecoverList))
			note.GET("/annex/download", ichat.HandlerFunc(handler.V1.ArticleAnnex.Download))
		}

		organize := v1.Group("/organize").Use(authorize)
		{
			organize.GET("/department/all", ichat.HandlerFunc(handler.V1.Organize.DepartmentList))
			organize.GET("/personnel/all", ichat.HandlerFunc(handler.V1.Organize.PersonnelList))
		}
	}

	// v2 接口
	v2 := router.Group("/api/v2")
	{
		v2.GET("/test", func(context *gin.Context) {
			context.JSON(200, map[string]any{"message": "success"})
		})
	}
}
