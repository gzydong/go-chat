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
func RegisterApiRoute(conf *config.Config, router *gin.Engine, handler *handler.Handler, session *cache.Session) {
	// 授权验证中间件
	authorize := jwtutil.Auth(conf.Jwt.Secret, "api", session)

	// v1 接口
	v1 := router.Group("/api/v1")
	{
		common := v1.Group("/common")
		{
			common.POST("/sms-code", handler.Api.Common.SmsCode)
			common.POST("/email-code", authorize, handler.Api.Common.EmailCode)
			common.GET("/setting", authorize, handler.Api.Common.Setting)
		}

		// 授权相关分组
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handler.Api.Auth.Login)                // 登录
			auth.POST("/register", handler.Api.Auth.Register)          // 注册
			auth.POST("/refresh", authorize, handler.Api.Auth.Refresh) // 刷新 Token
			auth.POST("/logout", authorize, handler.Api.Auth.Logout)   // 退出登录
			auth.POST("/forget", handler.Api.Auth.Forget)              // 找回密码
		}

		// 用户相关分组
		user := v1.Group("/users").Use(authorize)
		{
			user.GET("/detail", handler.Api.User.Detail)                   // 获取个人信息
			user.GET("/setting", handler.Api.User.Setting)                 // 获取个人信息
			user.POST("/change/detail", handler.Api.User.ChangeDetail)     // 修改用户信息
			user.POST("/change/password", handler.Api.User.ChangePassword) // 修改用户密码
			user.POST("/change/mobile", handler.Api.User.ChangeMobile)     // 修改用户手机号
			user.POST("/change/email", handler.Api.User.ChangeEmail)       // 修改用户邮箱
		}

		contact := v1.Group("/contact").Use(authorize)
		{
			contact.GET("/list", handler.Api.Contact.List)               // 联系人列表
			contact.GET("/search", handler.Api.Contact.Search)           // 搜索联系人
			contact.GET("/detail", handler.Api.Contact.Detail)           // 搜索联系人
			contact.POST("/delete", handler.Api.Contact.Delete)          // 删除联系人
			contact.POST("/edit-remark", handler.Api.Contact.EditRemark) // 编辑联系人备注

			// 联系人申请相关
			contact.GET("/apply/records", handler.Api.ContactsApply.List)              // 联系人申请列表
			contact.POST("/apply/create", handler.Api.ContactsApply.Create)            // 添加联系人申请
			contact.POST("/apply/accept", handler.Api.ContactsApply.Accept)            // 同意人申请列表
			contact.POST("/apply/decline", handler.Api.ContactsApply.Decline)          // 拒绝人申请列表
			contact.GET("/apply/unread-num", handler.Api.ContactsApply.ApplyUnreadNum) // 联系人申请未读数
		}

		// 聊天群相关分组
		userGroup := v1.Group("/group").Use(authorize)
		{
			userGroup.GET("/list", handler.Api.Group.GetGroups)            // 群组列表
			userGroup.GET("/overt/list", handler.Api.Group.OvertList)      // 公开群组列表
			userGroup.GET("/detail", handler.Api.Group.Detail)             // 群组详情
			userGroup.POST("/create", handler.Api.Group.Create)            // 创建群组
			userGroup.POST("/dismiss", handler.Api.Group.Dismiss)          // 解散群组
			userGroup.POST("/invite", handler.Api.Group.Invite)            // 邀请加入群组
			userGroup.POST("/secede", handler.Api.Group.SignOut)           // 退出群组
			userGroup.POST("/setting", handler.Api.Group.Setting)          // 设置群组信息
			userGroup.POST("/handover", handler.Api.Group.Handover)        // 群主转让
			userGroup.POST("/assign-admin", handler.Api.Group.AssignAdmin) // 分配管理员
			userGroup.POST("/no-speak", handler.Api.Group.NoSpeak)         // 修改禁言状态

			// 群成员相关
			userGroup.GET("/member/list", handler.Api.Group.GetMembers)          // 群成员列表
			userGroup.GET("/member/invites", handler.Api.Group.GetInviteFriends) // 群成员列表
			userGroup.POST("/member/remove", handler.Api.Group.RemoveMembers)    // 移出指定群成员
			userGroup.POST("/member/remark", handler.Api.Group.EditRemark)       // 设置群名片

			// 群公告相关
			userGroup.GET("/notice/list", handler.Api.GroupNotice.List)             // 群公告列表
			userGroup.POST("/notice/edit", handler.Api.GroupNotice.CreateAndUpdate) // 添加或编辑群公告
			userGroup.POST("/notice/delete", handler.Api.GroupNotice.Delete)        // 删除群公告

			// 群申请
			userGroup.POST("/apply/create", handler.Api.GroupApply.Create) // 提交入群申请
			userGroup.POST("/apply/delete", handler.Api.GroupApply.Delete) // 申请入群申请
			userGroup.POST("/apply/agree", handler.Api.GroupApply.Agree)   // 同意入群申请
			userGroup.GET("/apply/list", handler.Api.GroupApply.List)      // 入群申请列表
		}

		talk := v1.Group("/talk").Use(authorize)
		{
			talk.GET("/list", handler.Api.Talk.List)                                   // 会话列表
			talk.POST("/create", handler.Api.Talk.Create)                              // 创建会话
			talk.POST("/delete", handler.Api.Talk.Delete)                              // 删除会话
			talk.POST("/topping", handler.Api.Talk.Top)                                // 置顶会话
			talk.POST("/disturb", handler.Api.Talk.Disturb)                            // 会话免打扰
			talk.GET("/records", handler.Api.TalkRecords.GetRecords)                   // 会话面板记录
			talk.GET("/records/history", handler.Api.TalkRecords.SearchHistoryRecords) // 历史会话记录
			talk.GET("/records/forward", handler.Api.TalkRecords.GetForwardRecords)    // 会话转发记录
			talk.GET("/records/file/download", handler.Api.TalkRecords.Download)       // 会话转发记录
			talk.POST("/unread/clear", handler.Api.Talk.ClearUnreadMessage)            // 清除会话未读数
		}

		talkMsg := v1.Group("/talk/message").Use(authorize)
		{
			talkMsg.POST("/text", handler.Api.TalkMessage.Text)              // 发送文本消息
			talkMsg.POST("/code", handler.Api.TalkMessage.Code)              // 发送代码消息
			talkMsg.POST("/image", handler.Api.TalkMessage.Image)            // 发送图片消息
			talkMsg.POST("/file", handler.Api.TalkMessage.File)              // 发送文件消息
			talkMsg.POST("/emoticon", handler.Api.TalkMessage.Emoticon)      // 发送表情包消息
			talkMsg.POST("/forward", handler.Api.TalkMessage.Forward)        // 发送转发消息
			talkMsg.POST("/card", handler.Api.TalkMessage.Card)              // 发送用户名片
			talkMsg.POST("/location", handler.Api.TalkMessage.Location)      // 发送位置消息
			talkMsg.POST("/collect", handler.Api.TalkMessage.Collect)        // 收藏会话表情图片
			talkMsg.POST("/revoke", handler.Api.TalkMessage.Revoke)          // 撤销聊天消息
			talkMsg.POST("/delete", handler.Api.TalkMessage.Delete)          // 删除聊天消息
			talkMsg.POST("/vote", handler.Api.TalkMessage.Vote)              // 发送投票消息
			talkMsg.POST("/vote/handle", handler.Api.TalkMessage.HandleVote) // 投票消息处理
		}

		emoticon := v1.Group("/emoticon").Use(authorize)
		{
			emoticon.GET("/list", handler.Api.Emoticon.CollectList)                // 表情包列表
			emoticon.POST("/customize/create", handler.Api.Emoticon.Upload)        // 添加自定义表情
			emoticon.POST("/customize/delete", handler.Api.Emoticon.DeleteCollect) // 删除自定义表情

			// 系統表情包
			emoticon.GET("/system/list", handler.Api.Emoticon.SystemList)            // 系统表情包列表
			emoticon.POST("/system/install", handler.Api.Emoticon.SetSystemEmoticon) // 添加或移除系统表情包
		}

		upload := v1.Group("/upload").Use(authorize)
		{
			upload.POST("/avatar", handler.Api.Upload.Avatar)
			upload.POST("/multipart/initiate", handler.Api.Upload.InitiateMultipart)
			upload.POST("/multipart", handler.Api.Upload.MultipartUpload)
		}

		note := v1.Group("/note").Use(authorize)
		{
			// 文章相关
			note.GET("/article/list", handler.Api.Article.List)
			note.POST("/article/editor", handler.Api.Article.Edit)
			note.GET("/article/detail", handler.Api.Article.Detail)
			note.POST("/article/delete", handler.Api.Article.Delete)
			note.POST("/article/upload/image", handler.Api.Article.Upload)
			note.POST("/article/recover", handler.Api.Article.Recover)
			note.POST("/article/move", handler.Api.Article.Move)
			note.POST("/article/asterisk", handler.Api.Article.Asterisk)
			note.POST("/article/tag", handler.Api.Article.Tag)
			note.POST("/article/forever/delete", handler.Api.Article.ForeverDelete)

			// 文章分类
			note.GET("/class/list", handler.Api.ArticleClass.List)
			note.POST("/class/editor", handler.Api.ArticleClass.Edit)
			note.POST("/class/delete", handler.Api.ArticleClass.Delete)
			note.POST("/class/sort", handler.Api.ArticleClass.Sort)

			// 文章标签
			note.GET("/tag/list", handler.Api.ArticleTag.List)
			note.POST("/tag/editor", handler.Api.ArticleTag.Edit)
			note.POST("/tag/delete", handler.Api.ArticleTag.Delete)

			// 文章附件
			note.POST("/annex/upload", handler.Api.ArticleAnnex.Upload)
			note.POST("/annex/delete", handler.Api.ArticleAnnex.Delete)
			note.POST("/annex/recover", handler.Api.ArticleAnnex.Recover)
			note.POST("/annex/forever/delete", handler.Api.ArticleAnnex.ForeverDelete)
			note.GET("/annex/recover/list", handler.Api.ArticleAnnex.RecoverList)
			note.GET("/annex/download", handler.Api.ArticleAnnex.Download)
		}

		organize := v1.Group("/organize").Use(authorize)
		{
			organize.GET("/department/all", handler.Api.Organize.DepartmentList)
			organize.GET("/personnel/all", handler.Api.Organize.PersonnelList)
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
