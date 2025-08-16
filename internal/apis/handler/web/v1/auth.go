package v1

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/encrypt/aesutil"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/model"

	"go-chat/api/pb/queue/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Auth struct {
	Config              *config.Config
	Redis               *redis.Client
	JwtTokenStorage     *cache.JwtTokenStorage
	RedisLock           *cache.RedisLock
	RobotRepo           *repo.Robot
	OAuthUsersRepo      *repo.OAuthUsers
	UsersRepo           *repo.Users
	SmsService          service.ISmsService
	UserService         service.IUserService
	ArticleClassService service.IArticleClassService
	Rsa                 rsautil.IRsa
	OauthService        service.IOAuthService
	AesUtil             aesutil.IAesUtil
}

// Login 登录接口
func (c *Auth) Login(ctx *core.Context) error {
	in := &web.AuthLoginRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	user, err := c.UserService.Login(ctx.GetContext(), in.Mobile, string(password))
	if err != nil {
		return ctx.Error(err)
	}

	data := jsonutil.Marshal(queue.UserLoginRequest{
		UserId:   int32(user.Id),
		IpAddr:   ctx.Context.ClientIP(),
		Platform: in.Platform,
		Agent:    ctx.Context.GetHeader("user-agent"),
		LoginAt:  time.Now().Format(time.DateTime),
	})

	if err := c.Redis.Publish(ctx.GetContext(), entity.LoginTopic, data).Err(); err != nil {
		logger.ErrorWithFields(
			"投递登录消息异常", err,
			queue.UserLoginRequest{
				UserId:   int32(user.Id),
				IpAddr:   ctx.Context.ClientIP(),
				Platform: in.Platform,
				Agent:    ctx.Context.GetHeader("user-agent"),
				LoginAt:  time.Now().Format(time.DateTime),
			},
		)
	}

	authorize, err := c.authorize(user.Id)
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.AuthLoginResponse{
		Type:        authorize.Type,
		AccessToken: authorize.AccessToken,
		ExpiresIn:   authorize.ExpiresIn,
	})
}

// Register 注册接口
func (c *Auth) Register(ctx *core.Context) error {
	in := &web.AuthRegisterRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if !utils.IsMobile(in.Mobile) {
		return ctx.InvalidParams("Mobile 格式错误")
	}

	// 验证短信验证码是否正确
	if !c.SmsService.Verify(ctx.GetContext(), entity.SmsRegisterChannel, in.Mobile, in.SmsCode) {
		return entity.ErrSmsCodeError
	}

	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	if _, err := c.UserService.Register(ctx.GetContext(), &service.UserRegisterOpt{
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Password: string(password),
		Platform: in.Platform,
	}); err != nil {
		return ctx.Error(err)
	}

	c.SmsService.Delete(ctx.GetContext(), entity.SmsRegisterChannel, in.Mobile)

	return ctx.Success(&web.AuthRegisterResponse{})
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *core.Context) error {
	token := middleware.GetAuthToken(ctx.Context)

	claims, err := jwtutil.ParseWithClaims[entity.WebClaims]([]byte(c.Config.Jwt.Secret), token)
	if err == nil {
		if ex := claims.ExpiresAt.Unix() - time.Now().Unix(); ex > 0 {
			_ = c.JwtTokenStorage.SetBlackList(ctx.GetContext(), token, time.Duration(ex)*time.Second)
		}
	}

	return ctx.Success(nil)
}

// Forget 账号找回接口
func (c *Auth) Forget(ctx *core.Context) error {
	in := &web.AuthForgetRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if !utils.IsMobile(in.Mobile) {
		return ctx.InvalidParams("Mobile 格式错误")
	}

	// 验证短信验证码是否正确
	if !c.SmsService.Verify(ctx.GetContext(), entity.SmsForgetAccountChannel, in.Mobile, in.SmsCode) {
		return entity.ErrSmsCodeError
	}

	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	if _, err := c.UserService.Forget(ctx.GetContext(), &service.UserForgetOpt{
		Mobile:   in.Mobile,
		Password: string(password),
		SmsCode:  in.SmsCode,
	}); err != nil {
		return ctx.Error(err)
	}

	c.SmsService.Delete(ctx.GetContext(), entity.SmsForgetAccountChannel, in.Mobile)

	return ctx.Success(&web.AuthForgetResponse{})
}

// OAuth 第三方登录
func (c *Auth) OAuth(ctx *core.Context) error {
	in := &web.AuthOauthRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uri, err := c.OauthService.GetAuthURL(ctx.GetContext(), model.OAuthType(in.OauthType))
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.AuthOauthResponse{Uri: uri})
}

// OAuthLogin 第三方授权登录
func (c *Auth) OAuthLogin(ctx *core.Context) error {
	in := &web.AuthOauthLoginRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	oAuthInfo, err := c.OauthService.HandleCallback(ctx.GetContext(), model.OAuthType(in.OauthType), in.Code, in.State)
	if err != nil {
		return ctx.Error(err)
	}

	// 有会员信息直接返回登录信息
	if oAuthInfo.UserId > 0 {
		authorize, err := c.authorize(int(oAuthInfo.UserId))
		if err != nil {
			return ctx.Error(err)
		}

		return ctx.Success(&web.AuthOauthLoginResponse{
			IsAuthorize: "Y",
			Authorize:   authorize,
		})
	}

	ciphertext, err := c.AesUtil.Encrypt(jsonutil.Encode(BindTokenInfo{
		Id:        oAuthInfo.Id,
		Type:      string(oAuthInfo.OAuthType),
		Timestamp: time.Now().Unix(),
	}))

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.AuthOauthLoginResponse{
		IsAuthorize: "N",
		BindToken:   ciphertext,
	})
}

// OAuthBind 第三方授权登录绑定
func (c *Auth) OAuthBind(ctx *core.Context) error {
	in := &web.AuthOAuthBindRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	decrypt, err := c.AesUtil.Decrypt(in.BindToken)
	if err != nil {
		return ctx.InvalidParams("BindToken 解密异常")
	}

	var data = BindTokenInfo{}
	if err := jsonutil.Unmarshal(decrypt, &data); err != nil {
		return ctx.Error(err)
	}

	info, err := c.OAuthUsersRepo.FindById(ctx.GetContext(), data.Id)
	if err != nil {
		return ctx.Error(err)
	}

	if info.UserId != 0 {
		authorize, err := c.authorize(int(info.UserId))
		if err != nil {
			return ctx.Error(err)
		}

		return ctx.Success(&web.AuthOAuthBindResponse{
			Authorize: authorize,
		})
	}

	if !c.SmsService.Verify(ctx.GetContext(), entity.SmsOauthBindChannel, in.Mobile, in.SmsCode) {
		return entity.ErrSmsCodeError
	}

	userId, err := c.UserService.OauthBind(ctx.GetContext(), in.Mobile, info)
	if err != nil {
		return err
	}

	c.SmsService.Delete(ctx.GetContext(), entity.SmsOauthBindChannel, in.Mobile)

	authorize, err := c.authorize(userId)
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.AuthOAuthBindResponse{
		Authorize: authorize,
	})
}

// 生成 JWT Token
func (c *Auth) authorize(uid int) (*web.Authorize, error) {
	token, err := jwtutil.NewTokenWithClaims(
		[]byte(c.Config.Jwt.Secret), entity.WebClaims{
			UserId: int32(uid),
		},
		func(c *jwt.RegisteredClaims) {
			c.Issuer = entity.JwtIssuerWeb
		},
		jwtutil.WithTokenExpiresAt(time.Duration(c.Config.Jwt.ExpiresTime)*time.Second),
	)

	if err != nil {
		return nil, err
	}

	return &web.Authorize{
		AccessToken: token,
		ExpiresIn:   int32(c.Config.Jwt.ExpiresTime),
		Type:        "Bearer",
	}, nil
}

type BindTokenInfo struct {
	Id        int32  `json:"id"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}
