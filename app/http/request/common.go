package request

// SmsCodeRequest 发送短信验证码接口验证
type SmsCodeRequest struct {
	Mobile  string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"手机号"`
	Channel string `form:"channel" json:"channel" binding:"required,oneof=login register forget" label:"短信渠道"`
}
