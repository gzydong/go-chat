package v1

import "github.com/gin-gonic/gin"

type ContactsApply struct {
}

func NewContactsApplyHandler() *ContactsApply {
	return &ContactsApply{}
}

// ApplyUnreadNum 获取好友申请未读数
func (c *ContactsApply) ApplyUnreadNum(ctx *gin.Context) {

}

// Create 创建联系人申请
func (c *ContactsApply) Create(ctx *gin.Context) {

}

// Accept 同意联系人添加申请
func (c *ContactsApply) Accept(ctx *gin.Context) {

}

// Decline 拒绝联系人添加申请
func (c *ContactsApply) Decline(ctx *gin.Context) {

}

// List 获取联系人列表
func (c *ContactsApply) List(ctx *gin.Context) {

}
