package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/apis/handler/web"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/middleware"
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
			common.POST("/send-sms-code", core.HandlerFunc(handler.V1.Common.SmsCode))
			common.POST("/send-email-code", authorize, core.HandlerFunc(handler.V1.Common.EmailCode))
		}

		// 授权相关分组
		auth := v1.Group("/auth")
		{
			auth.POST("/login", core.HandlerFunc(handler.V1.Auth.Login))                // 登录
			auth.POST("/register", core.HandlerFunc(handler.V1.Auth.Register))          // 注册
			auth.POST("/refresh", authorize, core.HandlerFunc(handler.V1.Auth.Refresh)) // 刷新 Token
			auth.POST("/logout", authorize, core.HandlerFunc(handler.V1.Auth.Logout))   // 退出登录
			auth.POST("/forget", core.HandlerFunc(handler.V1.Auth.Forget))              // 找回密码
		}

		// 用户相关分组
		user := v1.Group("/user").Use(authorize)
		{
			user.GET("/detail", core.HandlerFunc(handler.V1.User.Detail))                   // 获取个人信息
			user.GET("/setting", core.HandlerFunc(handler.V1.User.Setting))                 // 获取个人信息
			user.POST("/update", core.HandlerFunc(handler.V1.User.ChangeDetail))            // 修改用户信息
			user.POST("/password/update", core.HandlerFunc(handler.V1.User.ChangePassword)) // 修改用户密码
			user.POST("/mobile/update", core.HandlerFunc(handler.V1.User.ChangeMobile))     // 修改用户手机号
			user.POST("/email/update", core.HandlerFunc(handler.V1.User.ChangeEmail))       // 修改用户邮箱
		}

		contact := v1.Group("/contact").Use(authorize)
		{
			contact.GET("/list", core.HandlerFunc(handler.V1.Contact.List))                   // 联系人列表
			contact.GET("/search", core.HandlerFunc(handler.V1.Contact.Search))               // 搜索联系人
			contact.GET("/detail", core.HandlerFunc(handler.V1.Contact.Detail))               // 搜索联系人
			contact.POST("/delete", core.HandlerFunc(handler.V1.Contact.Delete))              // 删除联系人
			contact.POST("/edit-remark", core.HandlerFunc(handler.V1.Contact.Remark))         // 编辑联系人备注
			contact.POST("/move-group", core.HandlerFunc(handler.V1.Contact.MoveGroup))       // 编辑联系人备注
			contact.POST("/online-status", core.HandlerFunc(handler.V1.Contact.OnlineStatus)) // 获取联系人在线状态

			// 联系人申请相关
			contact.GET("/apply/records", core.HandlerFunc(handler.V1.ContactApply.List))              // 联系人申请列表
			contact.POST("/apply/create", core.HandlerFunc(handler.V1.ContactApply.Create))            // 添加联系人申请
			contact.POST("/apply/accept", core.HandlerFunc(handler.V1.ContactApply.Accept))            // 同意人申请列表
			contact.POST("/apply/decline", core.HandlerFunc(handler.V1.ContactApply.Decline))          // 拒绝人申请列表
			contact.GET("/apply/unread-num", core.HandlerFunc(handler.V1.ContactApply.ApplyUnreadNum)) // 联系人申请未读数

			// 联系人分组
			contact.GET("/group/list", core.HandlerFunc(handler.V1.ContactGroup.List))    // 联系人分组列表
			contact.POST("/group/update", core.HandlerFunc(handler.V1.ContactGroup.Save)) // 联系人分组排序
		}

		// 聊天群相关分组
		userGroup := v1.Group("/group").Use(authorize)
		{
			userGroup.GET("/list", core.HandlerFunc(handler.V1.Group.List))                 // 群组列表
			userGroup.GET("/overt-list", core.HandlerFunc(handler.V1.Group.OvertList))      // 公开群组列表
			userGroup.GET("/detail", core.HandlerFunc(handler.V1.Group.Detail))             // 群组详情
			userGroup.POST("/create", core.HandlerFunc(handler.V1.Group.Create))            // 创建群组
			userGroup.POST("/dismiss", core.HandlerFunc(handler.V1.Group.Dismiss))          // 解散群组
			userGroup.POST("/invite", core.HandlerFunc(handler.V1.Group.Invite))            // 邀请加入群组
			userGroup.POST("/secede", core.HandlerFunc(handler.V1.Group.Secede))            // 退出群组
			userGroup.POST("/update", core.HandlerFunc(handler.V1.Group.Update))            // 设置群组信息
			userGroup.POST("/transfer", core.HandlerFunc(handler.V1.Group.Transfer))        // 群主转让
			userGroup.POST("/assign-admin", core.HandlerFunc(handler.V1.Group.AssignAdmin)) // 分配管理员
			userGroup.POST("/mute", core.HandlerFunc(handler.V1.Group.Mute))                // 修改群禁言状态
			userGroup.POST("/member-mute", core.HandlerFunc(handler.V1.Group.MemberMute))   // 修改群成员禁言状态
			userGroup.POST("/overt", core.HandlerFunc(handler.V1.Group.Overt))              // 修改群公开状态

			// 群投票相关
			userGroup.POST("/vote/create", core.HandlerFunc(handler.V1.GroupVote.Create)) // 创建群投票
			userGroup.POST("/vote/submit", core.HandlerFunc(handler.V1.GroupVote.Submit)) // 投票提交
			userGroup.POST("/vote/detail", core.HandlerFunc(handler.V1.GroupVote.Detail)) // 投票详情

			// 群成员相关
			userGroup.GET("/invite-list", core.HandlerFunc(handler.V1.Group.GetInviteFriends))             // 待邀请入群好友列表
			userGroup.GET("/member/list", core.HandlerFunc(handler.V1.Group.Members))                      // 群成员列表
			userGroup.POST("/member/remove", core.HandlerFunc(handler.V1.Group.RemoveMember))              // 移出指定群成员
			userGroup.POST("/member/update-remark", core.HandlerFunc(handler.V1.Group.UpdateMemberRemark)) // 设置群名片

			// 群公告相关
			userGroup.POST("/notice/edit", core.HandlerFunc(handler.V1.GroupNotice.CreateAndUpdate)) // 添加或编辑群公告

			// 群申请
			userGroup.POST("/apply/create", core.HandlerFunc(handler.V1.GroupApply.Create))        // 提交入群申请
			userGroup.POST("/apply/agree", core.HandlerFunc(handler.V1.GroupApply.Agree))          // 同意入群申请
			userGroup.POST("/apply/decline", core.HandlerFunc(handler.V1.GroupApply.Decline))      // 拒绝入群申请
			userGroup.GET("/apply/list", core.HandlerFunc(handler.V1.GroupApply.List))             // 入群申请列表
			userGroup.GET("/apply/all", core.HandlerFunc(handler.V1.GroupApply.All))               // 入群申请列表
			userGroup.GET("/apply/unread", core.HandlerFunc(handler.V1.GroupApply.ApplyUnreadNum)) // 入群申请未读
		}

		talk := v1.Group("/talk").Use(authorize)
		{
			talk.GET("/list", core.HandlerFunc(handler.V1.Talk.List))                                   // 会话列表
			talk.POST("/create", core.HandlerFunc(handler.V1.Talk.Create))                              // 创建会话
			talk.POST("/delete", core.HandlerFunc(handler.V1.Talk.Delete))                              // 删除会话
			talk.POST("/topping", core.HandlerFunc(handler.V1.Talk.Top))                                // 置顶会话
			talk.POST("/disturb", core.HandlerFunc(handler.V1.Talk.Disturb))                            // 会话免打扰
			talk.GET("/records", core.HandlerFunc(handler.V1.TalkRecords.GetRecords))                   // 会话面板记录
			talk.GET("/history-records", core.HandlerFunc(handler.V1.TalkRecords.SearchHistoryRecords)) // 历史会话记录
			talk.GET("/forward-records", core.HandlerFunc(handler.V1.TalkRecords.GetForwardRecords))    // 会话转发记录
			talk.GET("/file-download", core.HandlerFunc(handler.V1.TalkRecords.Download))               // 下载文件
			talk.POST("/clear-unread", core.HandlerFunc(handler.V1.Talk.ClearUnreadMessage))            // 清除会话未读数
		}

		talkMessage := v1.Group("/talk/message").Use(authorize)
		{
			talkMessage.POST("/send", core.HandlerFunc(handler.V1.Message.Send))         // 发送文本消息
			talkMessage.POST("/revoke", core.HandlerFunc(handler.V1.TalkMessage.Revoke)) // 撤销聊天消息
			talkMessage.POST("/delete", core.HandlerFunc(handler.V1.TalkMessage.Delete)) // 删除聊天消息
		}

		emoticon := v1.Group("/emoticon").Use(authorize)
		{
			emoticon.GET("/customize/list", core.HandlerFunc(handler.V1.Emoticon.List))      // 表情包列表
			emoticon.POST("/customize/upload", core.HandlerFunc(handler.V1.Emoticon.Upload)) // 上传自定义表情
			emoticon.POST("/customize/create", core.HandlerFunc(handler.V1.Emoticon.Create)) // 添加自定义表情
			emoticon.POST("/customize/delete", core.HandlerFunc(handler.V1.Emoticon.Delete)) // 删除自定义表情
		}

		upload := v1.Group("/upload").Use(authorize)
		{
			upload.POST("/media-file", core.HandlerFunc(handler.V1.Upload.Image))
			upload.POST("/init-multipart", core.HandlerFunc(handler.V1.Upload.InitiateMultipart))
			upload.POST("/multipart", core.HandlerFunc(handler.V1.Upload.MultipartUpload))
		}

		note := v1.Group("/note").Use(authorize)
		{
			// 文章相关
			note.GET("/article/list", core.HandlerFunc(handler.V1.Article.List))                     // 文章列表
			note.POST("/article/recycle-list", core.HandlerFunc(handler.V1.Article.RecycleList))     // 回收站文章列表
			note.POST("/article/editor", core.HandlerFunc(handler.V1.Article.Editor))                // 编辑文章
			note.GET("/article/detail", core.HandlerFunc(handler.V1.Article.Detail))                 // 文章详情
			note.POST("/article/delete", core.HandlerFunc(handler.V1.Article.Delete))                // 删除文章
			note.POST("/article/forever-delete", core.HandlerFunc(handler.V1.Article.ForeverDelete)) // 永久删除文章
			note.POST("/article/recover-delete", core.HandlerFunc(handler.V1.Article.Recover))       // 恢复已删除文章
			note.POST("/article/move-classify", core.HandlerFunc(handler.V1.Article.MoveClassify))   // 移动分类
			note.POST("/article/collect", core.HandlerFunc(handler.V1.Article.Collect))              // 收藏文章
			note.POST("/article/update-tag", core.HandlerFunc(handler.V1.Article.UpdateTag))         // 更新文章标签

			// 文章分类
			note.GET("/classify/list", core.HandlerFunc(handler.V1.ArticleClass.List))
			note.POST("/classify/create", core.HandlerFunc(handler.V1.ArticleClass.Edit))
			note.POST("/classify/update", core.HandlerFunc(handler.V1.ArticleClass.Edit))
			note.POST("/classify/delete", core.HandlerFunc(handler.V1.ArticleClass.Delete))
			note.POST("/classify/sort", core.HandlerFunc(handler.V1.ArticleClass.Sort))

			// 文章标签
			note.GET("/tag/list", core.HandlerFunc(handler.V1.ArticleTag.List))
			note.POST("/tag/create", core.HandlerFunc(handler.V1.ArticleTag.Edit))
			note.POST("/tag/update", core.HandlerFunc(handler.V1.ArticleTag.Edit))
			note.POST("/tag/delete", core.HandlerFunc(handler.V1.ArticleTag.Delete))

			// 文章附件
			note.POST("/annex/upload", core.HandlerFunc(handler.V1.ArticleAnnex.Upload))
			note.POST("/annex/delete", core.HandlerFunc(handler.V1.ArticleAnnex.Delete))
			note.POST("/annex/recover", core.HandlerFunc(handler.V1.ArticleAnnex.Recover))
			note.POST("/annex/forever-delete", core.HandlerFunc(handler.V1.ArticleAnnex.ForeverDelete))
			note.GET("/annex/recycle-list", core.HandlerFunc(handler.V1.ArticleAnnex.RecycleList))
			note.GET("/annex/download", core.HandlerFunc(handler.V1.ArticleAnnex.Download))
		}

		organize := v1.Group("/organize").Use(authorize)
		{
			organize.GET("/department/all", core.HandlerFunc(handler.V1.Organize.DepartmentList))
			organize.GET("/personnel/all", core.HandlerFunc(handler.V1.Organize.PersonnelList))
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
