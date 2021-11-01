package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/strutil"
	"go-chat/app/service"
	"strings"
)

type Talk struct {
	service         *service.TalkService
	talkListService *service.TalkListService
	redisLock       *cache.RedisLock
}

func NewTalkHandler(
	service *service.TalkService,
	talkListService *service.TalkListService,
	redisLock *cache.RedisLock,
) *Talk {
	return &Talk{service, talkListService, redisLock}
}

// List 会话列表
func (c *Talk) List(ctx *gin.Context) {

}

// Create 创建会话列表
func (c *Talk) Create(ctx *gin.Context) {
	params := &request.TalkListCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	agent := strings.TrimSpace(ctx.GetHeader("user-agent"))
	if agent != "" {
		agent = strutil.Md5([]byte(agent))
	}

	lockKey := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, params.ReceiverId, params.TalkType, agent)
	if !c.redisLock.Lock(ctx.Request.Context(), lockKey, 20) {
		response.BusinessError(ctx, "创建失败")
		return
	}

	_, err := c.talkListService.Create(ctx.Request.Context(), uid, params)

	if err != nil {
		fmt.Println(err.Error())
	}

	response.Success(ctx, gin.H{})
}

// Delete 删除列表
func (c *Talk) Delete(ctx *gin.Context) {

}

// Top 置顶列表
func (c *Talk) Top(ctx *gin.Context) {

}

// Disturb 会话免打扰
func (c *Talk) Disturb(ctx *gin.Context) {

}
