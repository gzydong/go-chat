package request

// ChangeDetailRequest ...
type ChangeDetailRequest struct {
	Avatar   string `form:"avatar" binding:"" label:"avatar"`
	Nickname string `form:"nickname" binding:"required,max=30" label:"nickname"`
	Gender   string `form:"gender" binding:"required,oneof=0 1 2" label:"gender"`
	Profile  string `form:"profile" binding:"max=255" label:"profile"`
}

// ChangePasswordRequest ...
type ChangePasswordRequest struct {
	OldPassword string `form:"old_password" binding:"required" label:"old_password"`
	NewPassword string `form:"new_password" binding:"required,min=6,max=16" label:"new_password"`
}

// ChangeMobileRequest ...
type ChangeMobileRequest struct {
	Mobile   string `form:"mobile" binding:"required,len=11,phone" label:"mobile"`
	Password string `form:"password" binding:"required" label:"password"`
	SmsCode  string `form:"sms_code" binding:"required,len=6,numeric" label:"sms_code"`
}

// ChangeEmailRequest ...
type ChangeEmailRequest struct {
	Email    string `form:"email" binding:"required" label:"email"`
	Password string `form:"password" binding:"required" label:"password"`
	Code     string `form:"code" binding:"required,len=6,numeric" label:"code"`
}
