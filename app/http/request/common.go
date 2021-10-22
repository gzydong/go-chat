package request

// SmsCodeRequest 发送短信验证码接口验证
type SmsCodeRequest struct {
	Mobile  string `form:"mobile" binding:"required,len=11,phone" label:"mobile"`
	Channel string `form:"channel" binding:"required,oneof=login register forget_account change_account" label:"channel"`
}

type EmailCodeRequest struct {
	Email string `form:"email" binding:"required" label:"email"`
}
