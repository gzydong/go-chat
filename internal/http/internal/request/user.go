package request

// ChangeDetailRequest ...
type ChangeDetailRequest struct {
	Avatar   string `form:"avatar" json:"avatar" binding:"" label:"avatar"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,max=30" label:"nickname"`
	Gender   int    `form:"gender" json:"gender" binding:"oneof=0 1 2" label:"gender"`
	Motto    string `form:"motto" json:"motto" binding:"max=255" label:"motto"`
}

// ChangePasswordRequest ...
type ChangePasswordRequest struct {
	OldPassword string `form:"old_password" json:"old_password" binding:"required" label:"old_password"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=6,max=16" label:"new_password"`
}

// ChangeMobileRequest ...
type ChangeMobileRequest struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"mobile"`
	Password string `form:"password" json:"password" binding:"required" label:"password"`
	SmsCode  string `form:"sms_code" json:"sms_code" binding:"required,len=6,numeric" label:"sms_code"`
}

// ChangeEmailRequest ...
type ChangeEmailRequest struct {
	Email    string `form:"email" json:"email" binding:"required" label:"email"`
	Password string `form:"password" json:"password" binding:"required" label:"password"`
	Code     string `form:"code" json:"code" binding:"required,len=6,numeric" label:"code"`
}
