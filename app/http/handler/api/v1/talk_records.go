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

	response.Success(ctx, records)
}

func (c *TalkRecords) SearchRecords(ctx *gin.Context) {

}
