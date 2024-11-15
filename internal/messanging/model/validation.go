package event

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func ValidateUUID() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		uuidToValidate := fl.Field().String()
		_, err := uuid.Parse(uuidToValidate)
		return err == nil
	}
}
