package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/api/handler"
	"go-chat/internal/pkg/ichat/middleware"
	"go-chat/internal/repository/cache"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handler *handler.Handler, session *cache.JwtTokenStorage) *gin.Engine {
	router := gin.New()

	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err any) {
		log.Println(err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"code": 500, "msg": "系统错误，请重试!!!"})
	}))
	router.Use(middleware.Cors(conf.Cors))
	// router.Use(middleware.AccessLog())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]any{"code": 200, "message": "hello world"})
	})

	router.GET("/health/check", func(c *gin.Context) {
		c.JSON(200, map[string]any{"status": "ok"})
	})

	RegisterWebRoute(conf.Jwt.Secret, router, handler.Api, session)
	RegisterAdminRoute(conf.Jwt.Secret, router, handler.Admin, session)
	RegisterOpenRoute(router, handler.Open)

	// 注册 debug 路由
	if conf.Debug() {
		RegisterDebugRoute(router)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, map[string]any{"code": 404, "message": "请求地址不存在"})
	})

	return router
}
