package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/pakg/im"
	"net/http"
)

type UserController struct {
}

// Detail 个人用户信息
func (u *UserController) Detail(c *gin.Context) {
	msg := c.DefaultQuery("message", "")

	im.Manager.DefaultChannel.SendMessage(&im.Message{
		Clients: make([]string, 0),
		IsAll:   true,
		Event:   "talk_type",
		Content: msg,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}

// ChangePassword 修改密码接口
func (u *UserController) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}

// ChangeMobile 修改手机号接口
func (u UserController) ChangeMobile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}

// ChangeEmail 修改邮箱接口
func (u UserController) ChangeEmail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}
