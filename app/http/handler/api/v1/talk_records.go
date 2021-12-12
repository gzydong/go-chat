package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/service"
)

type TalkRecords struct {
	service *service.TalkRecordsService
}

func NewTalkRecordsHandler(service *service.TalkRecordsService) *TalkRecords {
	return &TalkRecords{
		service: service,
	}
}

// GetRecords 获取会话记录
func (c *TalkRecords) GetRecords(ctx *gin.Context) {
	params := &request.TalkRecordsRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	records, err := c.service.GetTalkRecords(ctx, &service.QueryTalkRecordsOpts{
		TalkType:   params.TalkType,
		UserId:     auth.GetAuthUserID(ctx),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Id
	}

	response.Success(ctx, gin.H{
		"limit":     params.Limit,
		"record_id": rid,
		"rows":      records,
	})
}

// SearchHistoryRecords 查询下会话记录
func (c *TalkRecords) SearchHistoryRecords(ctx *gin.Context) {
	c.GetRecords(ctx)
}

// GetForwardRecords 获取转发记录
func (c *TalkRecords) GetForwardRecords(ctx *gin.Context) {
	params := &request.TalkForwardRecordsRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	records, err := c.service.GetForwardRecords(ctx.Request.Context(), int64(params.RecordId))
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"rows": records,
	})
}
