package request

// LoginRequest 登录接口验证
type LoginRequest struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required" label:"登录账号"`
	Password string `form:"password" json:"password" binding:"required" label:"登录密码"`
	Platform string `form:"platform" json:"platform" binding:"required,oneof=h5 ios windows mac web"`
}

// RegisterRequest 注册接口验证
type RegisterRequest struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=30" label:"账号昵称"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"手机号"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16" label:"密码"`
	SmsCode  string `form:"sms_code" json:"sms_code" binding:"required" label:"验证码"`
	Platform string `form:"platform" json:"platform" binding:"required,oneof=h5 ios windows mac web" label:"登录平台"`
}

// ForgetRequest 账号找回接口验证
type ForgetRequest struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"手机号"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16" label:"密码"`
	SmsCode  string `form:"sms_code" json:"sms_code" binding:"required,len=6" label:"验证码"`
}
