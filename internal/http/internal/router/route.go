package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/middleware"
	"go-chat/internal/http/internal/response"
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
		response.Success(c, "")
	})

	router.GET("/health/check", func(c *gin.Context) {
		c.JSON(200, entity.H{"status": "ok"})
	})

	RegisterApiRoute(conf, router, handler.Api, tokenCache)
	RegisterAdminRoute(conf, router, handler.Admin, tokenCache)
	RegisterOpenRoute(conf, router, handler.Open, tokenCache)

	router.NoRoute(func(c *gin.Context) {
		response.NewError(c, 404, "请求地址不存在")
	})

	return router
}
