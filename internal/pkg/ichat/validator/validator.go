package validator

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	_ = Initialize()
}

func Initialize() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		chinese := zh.New()
		uni := ut.New(chinese)
		trans, _ = uni.GetTranslator("zh")

		// 注册一个函数，获取struct tag里自定义的label作为字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})

		registerCustomValidator(v, trans)

		return zhTranslations.RegisterDefaultTranslations(v, trans)
	}

	return nil
}

func Translate(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, err := range errs {
			return err.Translate(trans)
		}
	}

	return err.Error()
}

func Validate(value interface{}) error {
	return binding.Validator.Engine().(*validator.Validate).Struct(value)
}

// registerCustomValidator 注册自定义验证器
func registerCustomValidator(v *validator.Validate, trans ut.Translator) {
	phone(v, trans)
	ids(v, trans)
}
