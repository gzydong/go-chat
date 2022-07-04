package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/middleware"
	"go-chat/internal/repository/cache"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handler *handler.Handler, tokenCache *cache.Session) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err interface{}) {

		// errorStr := fmt.Sprintf("[Recovery] %s panic recovered: %s\n%s", timeutil.FormatDatetime(time.Now()), err, string(debug.Stack(4)))
		//
		// fmt.Println(errorStr)

		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "系统错误，请重试!!!"})
	}))

	// 注册跨域中间件
	router.Use(middleware.Cors(conf))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, entity.H{"code": 200, "message": "hello world"})
	})

	router.GET("/health/check", func(c *gin.Context) {
		c.JSON(200, entity.H{"status": "ok"})
	})

	RegisterWebRoute(conf.Jwt.Secret, router, handler.Api, tokenCache)
	RegisterAdminRoute(conf.Jwt.Secret, router, handler.Admin, tokenCache)
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
