package util

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/hashicorp/go-version"
	"golang.org/x/text/language"
)

var NameExp = `^[_\w][-_\w]*$`

var NameReg = regexp.MustCompile(NameExp)

func ValidateVarName(name string) error {
	if !NameReg.MatchString(name) {
		return fmt.Errorf("input must be adhere to %s", NameExp)
	}
	return nil
}

func ValidateVersion(v string) error {
	if _, err := version.NewVersion(v); err != nil {
		return fmt.Errorf("invalid version format")
	}
	return nil
}

var (
	validate = validator.New()
	trans    ut.Translator
)

var matcher = language.NewMatcher([]language.Tag{
	language.English,
	language.Chinese,
})

func getLang() string {
	lang, _ := os.LookupEnv("LANG")
	tag, _ := language.MatchStrings(matcher, strings.Split(lang, ".")[0])
	b, _ := tag.Base()
	return b.String()
}

func registerZhTrans(trans ut.Translator) {
	validate.RegisterTranslation("version", trans, func(ut ut.Translator) error {
		return ut.Add("version", "{0}必须符合版本号格式", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("version", fe.Field())
		return t
	})
	validate.RegisterTranslation("var", trans, func(ut ut.Translator) error {
		return ut.Add("var", "{0}必须符合格式'{1}'", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("var", fe.Field(), NameExp)
		return t
	})
}

func registerEnTrans(trans ut.Translator) {
	validate.RegisterTranslation("version", trans, func(ut ut.Translator) error {
		return ut.Add("version", "{0} should be in version form", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("version", fe.Field())
		return t
	})
	validate.RegisterTranslation("var", trans, func(ut ut.Translator) error {
		return ut.Add("var", "{0} should be in form '{1}'", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("var", fe.Field(), NameExp)
		return t
	})
}

func init() {
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, en, zh)

	ent, _ := uni.GetTranslator("en")
	zht, _ := uni.GetTranslator("zh")
	en_translations.RegisterDefaultTranslations(validate, ent)
	zh_translations.RegisterDefaultTranslations(validate, zht)

	registerEnTrans(ent)
	registerZhTrans(zht)

	trans, _ = uni.GetTranslator(getLang())

	validate.RegisterValidation("version", func(fl validator.FieldLevel) bool {
		return ValidateVersion(fl.Field().String()) == nil
	})
	validate.RegisterValidation("var", func(fl validator.FieldLevel) bool {
		return ValidateVarName(fl.Field().String()) == nil
	})

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		if name := field.Tag.Get("flag"); name != "" {
			return name
		}
		if name := field.Tag.Get("alias"); name != "" {
			return name
		}
		return field.Name
	})
}
