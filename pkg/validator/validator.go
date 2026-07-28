package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func Validate(s interface{}) error {
	if validate == nil {
		Init()
	}
	return validate.Struct(s)
}

func TranslateErrors(err error) map[string]string {
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errors[field] = fmt.Sprintf("%s is required", field)
		case "email":
			errors[field] = fmt.Sprintf("%s must be a valid email", field)
		case "min":
			errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
		case "max":
			errors[field] = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
		case "uuid":
			errors[field] = fmt.Sprintf("%s must be a valid UUID", field)
		default:
			errors[field] = fmt.Sprintf("%s is invalid", field)
		}
	}
	return errors
}
