package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/handler"
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
			common.POST("/sms-code", handler.Common.SmsCode)
			common.POST("/email-code", authorize, handler.Common.EmailCode)
			common.GET("/setting", authorize, handler.Common.Setting)
		}

		// 授权相关分组
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handler.Auth.Login)                // 登录
			auth.POST("/register", handler.Auth.Register)          // 注册
			auth.POST("/refresh", authorize, handler.Auth.Refresh) // 刷新 Token
			auth.POST("/logout", authorize, handler.Auth.Logout)   // 退出登录
			auth.POST("/forget", handler.Auth.Forget)              // 找回密码
			auth.GET("/test", handler.Auth.Test)                   // 找回密码
		}

		// 用户相关分组
		user := v1.Group("/users").Use(authorize)
		{
			user.GET("/detail", handler.User.Detail)                   // 获取个人信息
			user.GET("/setting", handler.User.Setting)                 // 获取个人信息
			user.POST("/change/detail", handler.User.ChangeDetail)     // 修改用户信息
			user.POST("/change/password", handler.User.ChangePassword) // 修改用户密码
			user.POST("/change/mobile", handler.User.ChangeMobile)     // 修改用户手机号
			user.POST("/change/email", handler.User.ChangeEmail)       // 修改用户邮箱
		}

		contact := v1.Group("/contact").Use(authorize)
		{
			contact.GET("/list", handler.Contact.List)               // 联系人列表
			contact.GET("/search", handler.Contact.Search)           // 搜索联系人
			contact.GET("/detail", handler.Contact.Detail)           // 搜索联系人
			contact.POST("/delete", handler.Contact.Delete)          // 删除联系人
			contact.POST("/edit-remark", handler.Contact.EditRemark) // 编辑联系人备注

			// 联系人申请相关
			contact.GET("/apply/records", handler.ContactsApply.List)              // 联系人申请列表
			contact.POST("/apply/create", handler.ContactsApply.Create)            // 添加联系人申请
			contact.POST("/apply/accept", handler.ContactsApply.Accept)            // 同意人申请列表
			contact.POST("/apply/decline", handler.ContactsApply.Decline)          // 拒绝人申请列表
			contact.GET("/apply/unread-num", handler.ContactsApply.ApplyUnreadNum) // 联系人申请未读数
		}

		// 聊天群相关分组
		userGroup := v1.Group("/group").Use(authorize)
		{
			userGroup.GET("/list", handler.Group.GetGroups)            // 群组列表
			userGroup.GET("/overt/list", handler.Group.OvertList)      // 公开群组列表
			userGroup.GET("/detail", handler.Group.Detail)             // 群组详情
			userGroup.POST("/create", handler.Group.Create)            // 创建群组
			userGroup.POST("/dismiss", handler.Group.Dismiss)          // 解散群组
			userGroup.POST("/invite", handler.Group.Invite)            // 邀请加入群组
			userGroup.POST("/secede", handler.Group.SignOut)           // 退出群组
			userGroup.POST("/setting", handler.Group.Setting)          // 设置群组信息
			userGroup.POST("/handover", handler.Group.Handover)        // 群主转让
			userGroup.POST("/assign-admin", handler.Group.AssignAdmin) // 分配管理员
			userGroup.POST("/no-speak", handler.Group.NoSpeak)         // 修改禁言状态

			// 群成员相关
			userGroup.GET("/member/list", handler.Group.GetMembers)          // 群成员列表
			userGroup.GET("/member/invites", handler.Group.GetInviteFriends) // 群成员列表
			userGroup.POST("/member/remove", handler.Group.RemoveMembers)    // 移出指定群成员
			userGroup.POST("/member/remark", handler.Group.EditRemark)       // 设置群名片

			// 群公告相关
			userGroup.GET("/notice/list", handler.GroupNotice.List)             // 群公告列表
			userGroup.POST("/notice/edit", handler.GroupNotice.CreateAndUpdate) // 添加或编辑群公告
			userGroup.POST("/notice/delete", handler.GroupNotice.Delete)        // 删除群公告

			// 群申请
			userGroup.POST("/apply/create", handler.GroupApply.Create) // 提交入群申请
			userGroup.POST("/apply/delete", handler.GroupApply.Delete) // 申请入群申请
			userGroup.POST("/apply/agree", handler.GroupApply.Agree)   // 同意入群申请
			userGroup.GET("/apply/list", handler.GroupApply.List)      // 入群申请列表
		}

		talk := v1.Group("/talk").Use(authorize)
		{
			talk.GET("/list", handler.Talk.List)                                   // 会话列表
			talk.POST("/create", handler.Talk.Create)                              // 创建会话
			talk.POST("/delete", handler.Talk.Delete)                              // 删除会话
			talk.POST("/topping", handler.Talk.Top)                                // 置顶会话
			talk.POST("/disturb", handler.Talk.Disturb)                            // 会话免打扰
			talk.GET("/records", handler.TalkRecords.GetRecords)                   // 会话面板记录
			talk.GET("/records/history", handler.TalkRecords.SearchHistoryRecords) // 历史会话记录
			talk.GET("/records/forward", handler.TalkRecords.GetForwardRecords)    // 会话转发记录
			talk.GET("/records/file/download", handler.TalkRecords.Download)       // 会话转发记录
			talk.POST("/unread/clear", handler.Talk.ClearUnreadMessage)            // 清除会话未读数
		}

		talkMsg := v1.Group("/talk/message").Use(authorize)
		{
			talkMsg.POST("/text", handler.TalkMessage.Text)              // 发送文本消息
			talkMsg.POST("/code", handler.TalkMessage.Code)              // 发送代码消息
			talkMsg.POST("/image", handler.TalkMessage.Image)            // 发送图片消息
			talkMsg.POST("/file", handler.TalkMessage.File)              // 发送文件消息
			talkMsg.POST("/emoticon", handler.TalkMessage.Emoticon)      // 发送表情包消息
			talkMsg.POST("/forward", handler.TalkMessage.Forward)        // 发送转发消息
			talkMsg.POST("/card", handler.TalkMessage.Card)              // 发送用户名片
			talkMsg.POST("/location", handler.TalkMessage.Location)      // 发送位置消息
			talkMsg.POST("/collect", handler.TalkMessage.Collect)        // 收藏会话表情图片
			talkMsg.POST("/revoke", handler.TalkMessage.Revoke)          // 撤销聊天消息
			talkMsg.POST("/delete", handler.TalkMessage.Delete)          // 删除聊天消息
			talkMsg.POST("/vote", handler.TalkMessage.Vote)              // 发送投票消息
			talkMsg.POST("/vote/handle", handler.TalkMessage.HandleVote) // 投票消息处理
		}

		emoticon := v1.Group("/emoticon").Use(authorize)
		{
			emoticon.GET("/list", handler.Emoticon.CollectList)                // 表情包列表
			emoticon.POST("/customize/create", handler.Emoticon.Upload)        // 添加自定义表情
			emoticon.POST("/customize/delete", handler.Emoticon.DeleteCollect) // 删除自定义表情

			// 系統表情包
			emoticon.GET("/system/list", handler.Emoticon.SystemList)            // 系统表情包列表
			emoticon.POST("/system/install", handler.Emoticon.SetSystemEmoticon) // 添加或移除系统表情包
		}

		upload := v1.Group("/upload").Use(authorize)
		{
			upload.POST("/avatar", handler.Upload.Avatar)
			upload.POST("/multipart/initiate", handler.Upload.InitiateMultipart)
			upload.POST("/multipart", handler.Upload.MultipartUpload)
		}

		note := v1.Group("/note").Use(authorize)
		{
			// 文章相关
			note.GET("/article/list", handler.Article.List)
			note.POST("/article/editor", handler.Article.Edit)
			note.GET("/article/detail", handler.Article.Detail)
			note.POST("/article/delete", handler.Article.Delete)
			note.POST("/article/upload/image", handler.Article.Upload)
			note.POST("/article/recover", handler.Article.Recover)
			note.POST("/article/move", handler.Article.Move)
			note.POST("/article/asterisk", handler.Article.Asterisk)
			note.POST("/article/tag", handler.Article.Tag)
			note.POST("/article/forever/delete", handler.Article.ForeverDelete)

			// 文章分类
			note.GET("/class/list", handler.ArticleClass.List)
			note.POST("/class/editor", handler.ArticleClass.Edit)
			note.POST("/class/delete", handler.ArticleClass.Delete)
			note.POST("/class/sort", handler.ArticleClass.Sort)

			// 文章标签
			note.GET("/tag/list", handler.ArticleTag.List)
			note.POST("/tag/editor", handler.ArticleTag.Edit)
			note.POST("/tag/delete", handler.ArticleTag.Delete)

			// 文章附件
			note.POST("/annex/upload", handler.ArticleAnnex.Upload)
			note.POST("/annex/delete", handler.ArticleAnnex.Delete)
			note.POST("/annex/recover", handler.ArticleAnnex.Recover)
			note.POST("/annex/forever/delete", handler.ArticleAnnex.ForeverDelete)
			note.GET("/annex/recover/list", handler.ArticleAnnex.RecoverList)
			note.GET("/annex/download", handler.ArticleAnnex.Download)
		}

		organize := v1.Group("/organize").Use(authorize)
		{
			organize.GET("/department/all", handler.Organize.DepartmentList)
			organize.GET("/personnel/all", handler.Organize.PersonnelList)
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
