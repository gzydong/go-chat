package v1

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/encrypt/aesutil"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/model"

	"go-chat/api/pb/queue/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

var _ web.IAuthHandler = (*Auth)(nil)

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

func (a *Auth) Login(ctx context.Context, in *web.AuthLoginRequest) (*web.AuthLoginResponse, error) {
	password, err := a.Rsa.Decrypt(in.Password)
	if err != nil {
		return nil, err
	}

	user, err := a.UserService.Login(ctx, in.Mobile, string(password))
	if err != nil {
		return nil, err
	}

	ip := ""
	userAgent := ""

	data := jsonutil.Marshal(queue.UserLoginRequest{
		UserId:   int32(user.Id),
		IpAddr:   ip,
		Platform: in.Platform,
		Agent:    userAgent,
		LoginAt:  time.Now().Format(time.DateTime),
	})

	if err := a.Redis.Publish(ctx, entity.LoginTopic, data).Err(); err != nil {
		logger.ErrorWithFields(
			"投递登录消息异常", err,
			queue.UserLoginRequest{
				UserId:   int32(user.Id),
				IpAddr:   ip,
				Platform: in.Platform,
				Agent:    userAgent,
				LoginAt:  time.Now().Format(time.DateTime),
			},
		)
	}

	authorize, err := a.authorize(user.Id)
	if err != nil {
		return nil, err
	}

	return &web.AuthLoginResponse{
		Type:        authorize.Type,
		AccessToken: authorize.AccessToken,
		ExpiresIn:   authorize.ExpiresIn,
	}, nil
}

func (a *Auth) Register(ctx context.Context, in *web.AuthRegisterRequest) (*web.AuthRegisterResponse, error) {
	if !utils.IsMobile(in.Mobile) {
		return nil, errorx.New(400, "手机号格式不对")
	}

	// 验证短信验证码是否正确
	if !a.SmsService.Verify(ctx, entity.SmsRegisterChannel, in.Mobile, in.SmsCode) {
		return nil, entity.ErrSmsCodeError
	}

	password, err := a.Rsa.Decrypt(in.Password)
	if err != nil {
		return nil, err
	}

	user, err := a.UserService.Register(ctx, &service.UserRegisterOpt{
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Password: string(password),
		Platform: in.Platform,
	})

	if err != nil {
		return nil, err
	}

	a.SmsService.Delete(ctx, entity.SmsRegisterChannel, in.Mobile)

	authorize, err := a.authorize(user.Id)
	if err != nil {
		return nil, err
	}

	return &web.AuthRegisterResponse{
		Type:        authorize.Type,
		AccessToken: authorize.AccessToken,
		ExpiresIn:   authorize.ExpiresIn,
	}, nil
}

func (a *Auth) Forget(ctx context.Context, in *web.AuthForgetRequest) (*web.AuthForgetResponse, error) {
	if !utils.IsMobile(in.Mobile) {
		return nil, errorx.New(400, "手机号格式不对")
	}

	// 验证短信验证码是否正确
	if !a.SmsService.Verify(ctx, entity.SmsForgetAccountChannel, in.Mobile, in.SmsCode) {
		return nil, entity.ErrSmsCodeError
	}

	password, err := a.Rsa.Decrypt(in.Password)
	if err != nil {
		return nil, err
	}

	if _, err := a.UserService.Forget(ctx, &service.UserForgetOpt{
		Mobile:   in.Mobile,
		Password: string(password),
		SmsCode:  in.SmsCode,
	}); err != nil {
		return nil, err
	}

	a.SmsService.Delete(ctx, entity.SmsForgetAccountChannel, in.Mobile)

	return &web.AuthForgetResponse{}, nil
}

func (a *Auth) Oauth(ctx context.Context, in *web.AuthOauthRequest) (*web.AuthOauthResponse, error) {
	uri, err := a.OauthService.GetAuthURL(ctx, model.OAuthType(in.OauthType))
	if err != nil {
		return nil, err
	}

	return &web.AuthOauthResponse{Uri: uri}, nil
}

func (a *Auth) OauthBind(ctx context.Context, in *web.AuthOAuthBindRequest) (*web.AuthOAuthBindResponse, error) {
	decrypt, err := a.AesUtil.Decrypt(in.BindToken)
	if err != nil {
		return nil, err
	}

	var data = BindTokenInfo{}
	if err := jsonutil.Unmarshal(decrypt, &data); err != nil {
		return nil, err
	}

	info, err := a.OAuthUsersRepo.FindById(ctx, data.Id)
	if err != nil {
		return nil, err
	}

	if info.UserId != 0 {
		authorize, err := a.authorize(int(info.UserId))
		if err != nil {
			return nil, err
		}

		return &web.AuthOAuthBindResponse{
			Authorize: authorize,
		}, nil
	}

	if !a.SmsService.Verify(ctx, entity.SmsOauthBindChannel, in.Mobile, in.SmsCode) {
		return nil, entity.ErrSmsCodeError
	}

	userId, err := a.UserService.OauthBind(ctx, in.Mobile, info)
	if err != nil {
		return nil, err
	}

	a.SmsService.Delete(ctx, entity.SmsOauthBindChannel, in.Mobile)

	authorize, err := a.authorize(userId)
	if err != nil {
		return nil, err
	}

	return &web.AuthOAuthBindResponse{
		Authorize: authorize,
	}, nil
}

func (a *Auth) OauthLogin(ctx context.Context, in *web.AuthOauthLoginRequest) (*web.AuthOauthLoginResponse, error) {

	oAuthInfo, err := a.OauthService.HandleCallback(ctx, model.OAuthType(in.OauthType), in.Code, in.State)
	if err != nil {
		return nil, err
	}

	// 有会员信息直接返回登录信息
	if oAuthInfo.UserId > 0 {
		authorize, err := a.authorize(int(oAuthInfo.UserId))
		if err != nil {
			return nil, err
		}

		return &web.AuthOauthLoginResponse{
			IsAuthorize: "Y",
			Authorize:   authorize,
		}, nil
	}

	ciphertext, err := a.AesUtil.Encrypt(jsonutil.Encode(BindTokenInfo{
		Id:        oAuthInfo.Id,
		Type:      string(oAuthInfo.OAuthType),
		Timestamp: time.Now().Unix(),
	}))

	if err != nil {
		return nil, err
	}

	return &web.AuthOauthLoginResponse{
		IsAuthorize: "N",
		BindToken:   ciphertext,
	}, nil
}

// 生成 JWT Token
func (a *Auth) authorize(uid int) (*web.Authorize, error) {
	token, err := jwtutil.NewTokenWithClaims(
		[]byte(a.Config.Jwt.Secret), entity.WebClaims{
			UserId: int32(uid),
		},
		func(c *jwt.RegisteredClaims) {
			c.Issuer = entity.JwtIssuerWeb
		},
		jwtutil.WithTokenExpiresAt(time.Duration(a.Config.Jwt.ExpiresTime)*time.Second),
	)

	if err != nil {
		return nil, err
	}

	return &web.Authorize{
		AccessToken: token,
		ExpiresIn:   int32(a.Config.Jwt.ExpiresTime),
		Type:        "Bearer",
	}, nil
}

type BindTokenInfo struct {
	Id        int32  `json:"id"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}
