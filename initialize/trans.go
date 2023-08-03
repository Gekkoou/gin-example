package initialize

import (
	"fmt"
	"gin-example/global"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var uni *ut.UniversalTranslator

func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 修改返回字段key的格式
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni = ut.New(enT, zhT, enT)

		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) error", locale)
		}

		switch locale {
		case "zh":
			zhtranslations.RegisterDefaultTranslations(v, global.Trans)
		case "en":
			entranslations.RegisterDefaultTranslations(v, global.Trans)
		default:
			entranslations.RegisterDefaultTranslations(v, global.Trans)
		}
		return
	}
	return
}
