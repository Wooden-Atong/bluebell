package controller

import (
	"bluebell/models"
	"fmt"
	"reflect"
	"strings"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/gin-gonic/gin/binding"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// 自定义SignUpParam结构体校验函数
func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.ParamSignUp)

	if su.Password != su.RePassword {
		// 输出错误提示信息，最后一个参数就是传递的param
		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")
	}
}

// 去除提示信息中的结构体名称
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		//一般会自动存成ParamSignUp.Username，根据.切割字符串
		res[field[strings.Index(field, ".")+1:]] = err //🌟strings.Index("ParamSignUp.Username",".")获取字符串某个字符的下标
	}
	return res
}

// 定义一个全局翻译器T
var trans ut.Translator

// 初始化校验翻译器
func InitTrans(locale string) (err error) {
	// 🌟gin框架内置的有默认validator，这里修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		//为SignUpParam注册自定义校验方法
		v.RegisterStructValidation(SignUpParamStructLevelValidation, models.ParamSignUp{})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}
