package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/tidwall/sjson"
	"go-chat/config"
	"go-chat/internal/apis/handler"
	"go-chat/internal/pkg/ichat/middleware"
	"go-chat/internal/repository/cache"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handler *handler.Handler, session *cache.JwtTokenStorage) *gin.Engine {
	router := gin.New()

	accessFilterRule := middleware.NewAccessFilterRule()
	accessFilterRule.Exclude("/api/v1/talk/records")
	accessFilterRule.Exclude("/api/v1/talk/history")
	accessFilterRule.Exclude("/api/v1/talk/forward")
	accessFilterRule.Exclude("/api/v1/talk/publish")
	accessFilterRule.AddRule("/api/v1/auth/login", func(info *middleware.RequestInfo) {
		info.RequestBody, _ = sjson.Set(info.RequestBody, `password`, "过滤敏感字段")
	})

	router.Use(middleware.Cors(conf.Cors))
	router.Use(middleware.AccessLog(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/logs/access.log", conf.Log.Path), // 日志文件的位置
		MaxSize:    100,                                              // 文件最大尺寸（以MB为单位）
		MaxBackups: 3,                                                // 保留的最大旧文件数量
		MaxAge:     7,                                                // 保留旧文件的最大天数
		Compress:   true,                                             // 是否压缩/归档旧文件
		LocalTime:  true,                                             // 使用本地时间创建时间戳
	}, accessFilterRule))

	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"code": 500, "msg": "系统错误，请重试!!!"})
	}))

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
