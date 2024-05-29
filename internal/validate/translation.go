package validate

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func registerZhTrans(trans ut.Translator) {
	V.RegisterTranslation("version", trans, func(ut ut.Translator) error {
		return ut.Add("version", "{0}必须符合版本号格式", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("version", fe.Field())
		return t
	})
	V.RegisterTranslation("var", trans, func(ut ut.Translator) error {
		return ut.Add("nature", "{0}必须符合格式: {1}", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("nature", fe.Field(), variableNameExp)
		return t
	})
}

func registerEnTrans(trans ut.Translator) {
	V.RegisterTranslation("version", trans, func(ut ut.Translator) error {
		return ut.Add("version", "{0} should be in version form", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("version", fe.Field())
		return t
	})
	V.RegisterTranslation("var", trans, func(ut ut.Translator) error {
		return ut.Add("var", "{0} should be in form '{1}'", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("var", fe.Field(), variableNameExp)
		return t
	})
}
