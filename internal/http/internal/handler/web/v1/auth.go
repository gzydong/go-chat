package v1

import (
	"strconv"
	"time"

	"go-chat/api/pb/message/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service/note"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/service"
)

type Auth struct {
	config             *config.Config
	userService        *service.UserService
	smsService         *service.SmsService
	session            *cache.JwtTokenStorage
	redisLock          *cache.RedisLock
	talkMessageService *service.TalkMessageService
	ipAddressService   *service.IpAddressService
	talkSessionService *service.TalkSessionService
	noteClassService   *note.ArticleClassService
	robotRepo          *repo.Robot
	message            *service.MessageService
}

func NewAuth(config *config.Config, userService *service.UserService, smsService *service.SmsService, session *cache.JwtTokenStorage, redisLock *cache.RedisLock, talkMessageService *service.TalkMessageService, ipAddressService *service.IpAddressService, talkSessionService *service.TalkSessionService, noteClassService *note.ArticleClassService, robotDao *repo.Robot, message *service.MessageService) *Auth {
	return &Auth{config: config, userService: userService, smsService: smsService, session: session, redisLock: redisLock, talkMessageService: talkMessageService, ipAddressService: ipAddressService, talkSessionService: talkSessionService, noteClassService: noteClassService, robotRepo: robotDao, message: message}
}

// Login 登录接口
func (c *Auth) Login(ctx *ichat.Context) error {

	params := &web.AuthLoginRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := c.userService.Login(params.Mobile, params.Password)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	// TODO 这里需需要异步处理
	root, _ := c.robotRepo.GetLoginRobot(ctx.Ctx())
	if root != nil {
		ip := ctx.Context.ClientIP()

		address, _ := c.ipAddressService.FindAddress(ip)

		_, _ = c.talkSessionService.Create(ctx.Ctx(), &service.TalkSessionCreateOpt{
			UserId:     user.Id,
			TalkType:   entity.ChatPrivateMode,
			ReceiverId: root.UserId,
			IsBoot:     true,
		})

		// 推送登录消息
		_ = c.message.SendLogin(ctx.Ctx(), user.Id, &message.LoginMessageRequest{
			Ip:       ip,
			Address:  address,
			Platform: params.Platform,
			Agent:    ctx.Context.GetHeader("user-agent"),
			Reason:   "常用设备登录",
		})
	}

	return ctx.Success(&web.AuthLoginResponse{
		Type:        "Bearer",
		AccessToken: c.token(user.Id),
		ExpiresIn:   int32(c.config.Jwt.ExpiresTime),
	})
}

// Register 注册接口
func (c *Auth) Register(ctx *ichat.Context) error {

	params := &web.AuthRegisterRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.smsService.Check(ctx.Ctx(), entity.SmsRegisterChannel, params.Mobile, params.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	if _, err := c.userService.Register(&service.UserRegisterOpt{
		Nickname: params.Nickname,
		Mobile:   params.Mobile,
		Password: params.Password,
		Platform: params.Platform,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	c.smsService.Delete(ctx.Ctx(), entity.SmsRegisterChannel, params.Mobile)

	return ctx.Success(&web.AuthRegisterResponse{})
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *ichat.Context) error {

	c.toBlackList(ctx)

	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *ichat.Context) error {

	c.toBlackList(ctx)

	return ctx.Success(&web.AuthRefreshResponse{
		Type:        "Bearer",
		AccessToken: c.token(ctx.UserId()),
		ExpiresIn:   int32(c.config.Jwt.ExpiresTime),
	})
}

// Forget 账号找回接口
func (c *Auth) Forget(ctx *ichat.Context) error {

	params := &web.AuthForgetRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.smsService.Check(ctx.Ctx(), entity.SmsForgetAccountChannel, params.Mobile, params.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	if _, err := c.userService.Forget(&service.UserForgetOpt{
		Mobile:   params.Mobile,
		Password: params.Password,
		SmsCode:  params.SmsCode,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	c.smsService.Delete(ctx.Ctx(), entity.SmsForgetAccountChannel, params.Mobile)

	return ctx.Success(&web.AuthForgetResponse{})
}

func (c *Auth) token(uid int) string {

	expiresAt := time.Now().Add(time.Second * time.Duration(c.config.Jwt.ExpiresTime))

	// 生成登录凭证
	token := jwt.GenerateToken("api", c.config.Jwt.Secret, &jwt.Options{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        strconv.Itoa(uid),
		Issuer:    "im.web",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token
}

// 设置黑名单
func (c *Auth) toBlackList(ctx *ichat.Context) {

	session := ctx.JwtSession()
	if session != nil {
		if ex := session.ExpiresAt - time.Now().Unix(); ex > 0 {
			_ = c.session.SetBlackList(ctx.Ctx(), session.Token, time.Duration(ex)*time.Second)
		}
	}
}
