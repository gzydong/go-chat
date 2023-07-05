package validator

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// phone 手机号验证器
func phone(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.MatchString("^1[3456789][0-9]{9}$", fl.Field().String())
		return matched
	})

	// 根据提供的标记注册翻译
	_ = v.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "手机号格式错误!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field(), fe.Field())
		return t
	})
}

// ids 逗号拼接ID，字符串验证
func ids(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterValidation("ids", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if value == "" {
			return true
		}

		matched, _ := regexp.MatchString("^\\d+(\\,\\d+)*$", value)
		return matched
	})

	_ = v.RegisterTranslation("ids", trans, func(ut ut.Translator) error {
		return ut.Add("ids", "ids 格式错误!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ids", fe.Field(), fe.Field())
		return t
	})
}
