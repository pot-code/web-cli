package validate

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

var naturalNameExp = `^[-_\w][-_\w\d]+$`
var natureNameReg = regexp.MustCompile(naturalNameExp)

func ValidateNatureName(name string) error {
	if !natureNameReg.MatchString(name) {
		return fmt.Errorf("input must be in form: %s", naturalNameExp)
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
	V = validator.New()
	T ut.Translator
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
	V.RegisterTranslation("version", trans, func(ut ut.Translator) error {
		return ut.Add("version", "{0}必须符合版本号格式", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("version", fe.Field())
		return t
	})
	V.RegisterTranslation("nature", trans, func(ut ut.Translator) error {
		return ut.Add("nature", "{0}必须符合格式: {1}", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("nature", fe.Field(), naturalNameExp)
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
	V.RegisterTranslation("nature", trans, func(ut ut.Translator) error {
		return ut.Add("var", "{0} should be in form '{1}'", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("var", fe.Field(), naturalNameExp)
		return t
	})
}

func init() {
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, en, zh)

	ent, _ := uni.GetTranslator("en")
	zht, _ := uni.GetTranslator("zh")
	en_translations.RegisterDefaultTranslations(V, ent)
	zh_translations.RegisterDefaultTranslations(V, zht)

	registerEnTrans(ent)
	registerZhTrans(zht)

	T, _ = uni.GetTranslator(getLang())

	V.RegisterValidation("version", func(fl validator.FieldLevel) bool {
		return ValidateVersion(fl.Field().String()) == nil
	})

	V.RegisterValidation("nature", func(fl validator.FieldLevel) bool {
		return ValidateNatureName(fl.Field().String()) == nil
	})

	V.RegisterTagNameFunc(func(field reflect.StructField) string {
		if name := field.Tag.Get("flag"); name != "" {
			return name
		}
		if name := field.Tag.Get("alias"); name != "" {
			return name
		}
		return field.Name
	})
}
