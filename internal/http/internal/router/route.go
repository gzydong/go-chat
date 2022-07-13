package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/middleware"
	"go-chat/internal/repository/cache"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handler *handler.Handler, session *cache.SessionStorage) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(middleware.Cors(conf))
	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err interface{}) {

		fmt.Println(err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.H{"code": 500, "msg": "系统错误，请重试!!!"})
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, entity.H{"code": 200, "message": "hello world"})
	})

	router.GET("/health/check", func(c *gin.Context) {
		c.JSON(200, entity.H{"status": "ok"})
	})

	RegisterWebRoute(conf.Jwt.Secret, router, handler.Api, session)
	RegisterAdminRoute(conf.Jwt.Secret, router, handler.Admin, session)
	RegisterOpenRoute(router, handler.Open)

	// 注册 debug 路由
	if conf.Debug() {
		RegisterDebugRoute(router)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, entity.H{"code": 404, "message": "请求地址不存在"})
	})

	return router
}
