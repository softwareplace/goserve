package utils

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

type RTranslation struct {
	Tag        string
	Message    string
	Validation func(validator.FieldError) []string
}

type RValidation struct {
	Tag       string
	Validator func(fl validator.FieldLevel) bool
}

type ValidationSetting struct {
	Language      string
	SpecSeparator string
	Translations  []RTranslation
	Validators    []RValidation
}

func DefaultPasswordValidations() RValidation {
	return RValidation{
		Tag: "password",
		Validator: func(fl validator.FieldLevel) bool {
			password := fl.Field().String()
			if len(password) < 8 {
				return false
			}
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
			hasSpecial := regexp.MustCompile(`[!@#$%^&*()\-_=+{}\[\]|;:'",.<>/?]`).MatchString(password)
			return hasLower && hasUpper && hasNumber && hasSpecial
		},
	}
}

func DefaultTranslations() []RTranslation {
	return []RTranslation{
		{
			Tag:     "required",
			Message: "{0} is a required field",
			Validation: func(fe validator.FieldError) []string {
				return []string{fe.Field()}
			},
		},
		{
			Tag:     "min",
			Message: "{0} must be at least {1} characters",
			Validation: func(fe validator.FieldError) []string {
				return []string{fe.Field(), fe.Param()}
			},
		},
		{
			Tag:     "max",
			Message: "{0} must be at most {1} characters",
			Validation: func(fe validator.FieldError) []string {
				return []string{fe.Field(), fe.Param()}
			},
		},
		{
			Tag:     "email",
			Message: "{0} must be a valid email address",
			Validation: func(fe validator.FieldError) []string {
				return []string{fe.Field()}
			},
		},
		{
			Tag:     "password",
			Message: "{0} must contain at least: 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 special character",
			Validation: func(fe validator.FieldError) []string {
				return []string{fe.Field()}
			},
		},
	}
}

func Default() *ValidationSetting {
	return &ValidationSetting{
		Language:      "en",
		SpecSeparator: "\n",
		Translations:  DefaultTranslations(),
		Validators: []RValidation{
			DefaultPasswordValidations(),
		},
	}
}

func StructValidation(target interface{}, setting ...ValidationSetting) error {
	var finalSetting ValidationSetting

	if len(setting) > 0 {
		finalSetting = setting[0]
	} else {
		finalSetting = *Default()
	}

	validate := validator.New()

	lang := en.New()
	uni := ut.New(lang, lang)
	trans, _ := uni.GetTranslator(finalSetting.Language)

	registerTranslation := func(tag string, message string, customFunc func(validator.FieldError) []string) {
		_ = validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, customFunc(fe)...)
			return t
		})
	}

	for _, translation := range finalSetting.Translations {
		registerTranslation(translation.Tag, translation.Message, translation.Validation)
	}

	for _, vr := range finalSetting.Validators {
		_ = validate.RegisterValidation(vr.Tag, vr.Validator)
	}

	err := validate.Struct(target)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errorMessages []string
			for _, e := range validationErrors {
				errorMessages = append(errorMessages, e.Translate(trans))
			}
			return errors.New(strings.Join(errorMessages, finalSetting.SpecSeparator))
		}
		return err
	}

	return nil
}
