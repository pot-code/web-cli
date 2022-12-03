package validate

import (
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"golang.org/x/text/language"
)

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
