package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequest(req interface{}) error {
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return err
		}
		return err
	}
	return nil
}
