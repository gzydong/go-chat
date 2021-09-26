package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"regexp"
)

var trans ut.Translator

func InitValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		chinese := zh.New()
		uni := ut.New(chinese)
		trans, _ = uni.GetTranslator("zh")

		// 注册一个函数，获取struct tag里自定义的label作为字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})

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

		return zhTranslations.RegisterDefaultTranslations(v, trans)
	}

	return nil
}

func Translate(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			return err.Translate(trans)
		}
	}

	return err.Error()
}

func init() {
	_ = InitValidator()
}
