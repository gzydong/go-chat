package request

// SmsCodeRequest 发送短信验证码接口验证
type SmsCodeRequest struct {
	Mobile  string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"mobile"`
	Channel string `form:"channel" json:"channel" binding:"required,oneof=login register forget_account change_account" label:"channel"`
}

type EmailCodeRequest struct {
	Email string `form:"email" json:"email" binding:"required" label:"email"`
}
