package v1

import (
	"strconv"
	"time"

	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/dao"
	"go-chat/internal/service/note"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/service"
)

type Auth struct {
	config             *config.Config
	userService        *service.UserService
	smsService         *service.SmsService
	session            *cache.SessionStorage
	redisLock          *cache.RedisLock
	talkMessageService *service.TalkMessageService
	ipAddressService   *service.IpAddressService
	talkSessionService *service.TalkSessionService
	noteClassService   *note.ArticleClassService
	robotDao           *dao.RobotDao
}

func NewAuth(config *config.Config, userService *service.UserService, smsService *service.SmsService, session *cache.SessionStorage, redisLock *cache.RedisLock, talkMessageService *service.TalkMessageService, ipAddressService *service.IpAddressService, talkSessionService *service.TalkSessionService, noteClassService *note.ArticleClassService, robotDao *dao.RobotDao) *Auth {
	return &Auth{config: config, userService: userService, smsService: smsService, session: session, redisLock: redisLock, talkMessageService: talkMessageService, ipAddressService: ipAddressService, talkSessionService: talkSessionService, noteClassService: noteClassService, robotDao: robotDao}
}

// Login 登录接口
func (c *Auth) Login(ctx *ichat.Context) error {

	params := &web.AuthLoginRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := c.userService.Login(params.Mobile, params.Password)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	root, _ := c.robotDao.FindLoginRobot()
	if root != nil {
		ip := ctx.Context.ClientIP()

		address, _ := c.ipAddressService.FindAddress(ip)

		_, _ = c.talkSessionService.Create(ctx.RequestContext(), &service.TalkSessionCreateOpt{
			UserId:     user.Id,
			TalkType:   entity.ChatPrivateMode,
			ReceiverId: root.UserId,
			IsBoot:     true,
		})

		// 推送登录消息
		_ = c.talkMessageService.SendLoginMessage(ctx.RequestContext(), &service.LoginMessageOpt{
			UserId:   user.Id,
			Ip:       ip,
			Address:  address,
			Platform: params.Platform,
			Agent:    ctx.Context.GetHeader("user-agent"),
		})
	}

	return ctx.Success(&web.AuthLoginResponse{
		Type:        "Bearer",
		AccessToken: c.token(user.Id),
		ExpiresIn:   int(c.config.Jwt.ExpiresTime),
	})
}

// Register 注册接口
func (c *Auth) Register(ctx *ichat.Context) error {

	params := &web.RegisterRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.smsService.CheckSmsCode(ctx.RequestContext(), entity.SmsRegisterChannel, params.Mobile, params.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	_, err := c.userService.Register(&service.UserRegisterOpt{
		Nickname: params.Nickname,
		Mobile:   params.Mobile,
		Password: params.Password,
		Platform: params.Platform,
	})
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	c.smsService.DeleteSmsCode(ctx.RequestContext(), entity.SmsRegisterChannel, params.Mobile)

	return ctx.Success(nil)
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
		ExpiresIn:   int(c.config.Jwt.ExpiresTime),
	})
}

// Forget 账号找回接口
func (c *Auth) Forget(ctx *ichat.Context) error {

	params := &web.ForgetRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.smsService.CheckSmsCode(ctx.RequestContext(), entity.SmsForgetAccountChannel, params.Mobile, params.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	if _, err := c.userService.Forget(&service.UserForgetOpt{
		Mobile:   params.Mobile,
		Password: params.Password,
		SmsCode:  params.SmsCode,
	}); err != nil {
		return ctx.BusinessError(err.Error())
	}

	c.smsService.DeleteSmsCode(ctx.RequestContext(), entity.SmsForgetAccountChannel, params.Mobile)

	return ctx.Success(nil)
}

func (c *Auth) token(uid int) string {

	expiresAt := time.Now().Add(time.Second * time.Duration(c.config.Jwt.ExpiresTime)).Unix()

	// 生成登录凭证
	token := jwt.GenerateToken("api", c.config.Jwt.Secret, &jwt.Options{
		ExpiresAt: expiresAt,
		Id:        strconv.Itoa(uid),
	})

	return token
}

// 设置黑名单
func (c *Auth) toBlackList(ctx *ichat.Context) {

	session := ctx.JwtSession()
	if session != nil {
		ex := session.ExpiresAt - time.Now().Unix()

		// 将 session 加入黑名单
		_ = c.session.SetBlackList(ctx.RequestContext(), session.Token, int(ex))
	}
}
