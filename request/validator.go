package request

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
)

func StructValidation(target interface{}) error {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")

	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return getErrorMessage(fld)
	})

	err := validate.Struct(target)

	if err != nil {
		errorMessage := strings.ReplaceAll(err.Error(), "'", "") + ""
		re := regexp.MustCompile(` .*\.`)
		errorMessage = re.ReplaceAllString(errorMessage, " ")
		return errors.New(errorMessage)
	}

	return nil
}

func getErrorMessage(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("error_message"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}
