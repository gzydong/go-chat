package v1

import (
	"fmt"
	"strconv"
	"time"

	"go-chat/api/web/v1"
	"go-chat/internal/http/internal/dto"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/service/note"

	"github.com/gin-gonic/gin"

	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/service"
)

type Auth struct {
	config             *config.Config
	userService        *service.UserService
	smsService         *service.SmsService
	session            *cache.Session
	redisLock          *cache.RedisLock
	talkMessageService *service.TalkMessageService
	ipAddressService   *service.IpAddressService
	talkSessionService *service.TalkSessionService
	noteClassService   *note.ArticleClassService
}

func NewAuthHandler(
	config *config.Config,
	userService *service.UserService,
	smsService *service.SmsService,
	session *cache.Session,
	redisLock *cache.RedisLock,
	talkMessageService *service.TalkMessageService,
	ipAddressService *service.IpAddressService,
	talkSessionService *service.TalkSessionService,
	noteClassService *note.ArticleClassService,
) *Auth {
	return &Auth{
		config:             config,
		userService:        userService,
		smsService:         smsService,
		session:            session,
		redisLock:          redisLock,
		talkMessageService: talkMessageService,
		ipAddressService:   ipAddressService,
		talkSessionService: talkSessionService,
		noteClassService:   noteClassService,
	}
}

// Login 登录接口
func (c *Auth) Login(ctx *gin.Context) {

	params := &request.LoginRequest{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	user, err := c.userService.Login(params.Mobile, params.Password)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	ip := ctx.ClientIP()

	address, _ := c.ipAddressService.FindAddress(ip)

	_, _ = c.talkSessionService.Create(ctx.Request.Context(), &service.TalkSessionCreateOpts{
		UserId:     user.Id,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: 4257,
		IsBoot:     true,
	})

	// 推送登录消息
	_ = c.talkMessageService.SendLoginMessage(ctx.Request.Context(), &service.LoginMessageOpts{
		UserId:   user.Id,
		Ip:       ip,
		Address:  address,
		Platform: params.Platform,
		Agent:    ctx.GetHeader("user-agent"),
	})

	response.Success(ctx, c.createToken(user.Id))
}

// Register 注册接口
func (c *Auth) Register(ctx *gin.Context) {
	params := &request.RegisterRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 验证短信验证码是否正确
	if !c.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsRegisterChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(ctx, "短信验证码填写错误！")
		return
	}

	_, err := c.userService.Register(&service.UserRegisterOpts{
		Nickname: params.Nickname,
		Mobile:   params.Mobile,
		Password: params.Password,
		Platform: params.Platform,
	})

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	c.smsService.DeleteSmsCode(ctx.Request.Context(), entity.SmsRegisterChannel, params.Mobile)

	response.Success(ctx, nil, "注册成功！")
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *gin.Context) {
	c.toBlackList(ctx)

	response.Success(ctx, nil, "已退出登录！")
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *gin.Context) {
	tokenInfo := c.createToken(jwtutil.GetUid(ctx))

	c.toBlackList(ctx)

	response.Success(ctx, tokenInfo)
}

// Forget 账号找回接口
func (c *Auth) Forget(ctx *gin.Context) {
	params := &request.ForgetRequest{}

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 验证短信验证码是否正确
	if !c.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(ctx, "短信验证码填写错误！")
		return
	}

	if _, err := c.userService.Forget(&service.UserForgetOpts{
		Mobile:   params.Mobile,
		Password: params.Password,
		SmsCode:  params.SmsCode,
	}); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	c.smsService.DeleteSmsCode(ctx.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile)

	response.Success(ctx, nil, "账号成功找回")
}

func (c *Auth) createToken(uid int) *dto.Token {

	expiresAt := time.Now().Add(time.Second * time.Duration(c.config.Jwt.ExpiresTime)).Unix()

	// 生成登录凭证
	token := jwtutil.GenerateToken("api", c.config.Jwt.Secret, &jwtutil.Options{
		ExpiresAt: expiresAt,
		Id:        strconv.Itoa(uid),
	})

	return &dto.Token{
		Type:      "Bearer",
		Token:     token,
		ExpiresIn: c.config.Jwt.ExpiresTime,
	}
}

// 设置黑名单
func (c *Auth) toBlackList(ctx *gin.Context) {
	info := ctx.GetStringMapString("jwt")

	expiresAt, _ := strconv.Atoi(info["expires_at"])

	ex := expiresAt - int(time.Now().Unix())

	// 将 session 加入黑名单
	_ = c.session.SetBlackList(ctx.Request.Context(), info["session"], ex)
}

func (c *Auth) Test(ctx *gin.Context) {
	params := &web.AuthLoginRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	fmt.Println("ansjkanskja", params)
}
