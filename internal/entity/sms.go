package entity

type SmsSendChannel string

const (
	SmsLoginChannel         = "login"
	SmsRegisterChannel      = "register"
	SmsForgetAccountChannel = "forget_account"
	SmsChangeAccountChannel = "change_account"
)
