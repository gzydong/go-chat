package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"net/http"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config) *gin.Engine {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok": "success",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "请求地址不存在",
		})
	})

	return router
}
