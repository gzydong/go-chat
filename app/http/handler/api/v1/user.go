package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/app/pkg/im"
	"go-chat/connect"
)

type User struct {
	MySQl *connect.MySQL
}

// Detail 个人用户信息
func (u *User) Detail(c *gin.Context) {
	msg := c.DefaultQuery("message", "")

	im.Manager.DefaultChannel.SendMessage(&im.SendMessage{
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
func (u *User) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}

// ChangeMobile 修改手机号接口
func (u User) ChangeMobile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}

// ChangeEmail 修改邮箱接口
func (u User) ChangeEmail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}
