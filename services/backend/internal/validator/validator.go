package validator

import (
	"strings"

	playground "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *playground.Validate
}

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New() *Validator {
	return &Validator{validate: playground.New()}
}

func (v *Validator) Struct(value any) []FieldError {
	if err := v.validate.Struct(value); err != nil {
		errs := make([]FieldError, 0)
		for _, fieldErr := range err.(playground.ValidationErrors) {
			field := strings.ToLower(fieldErr.Field()[:1]) + fieldErr.Field()[1:]
			errs = append(errs, FieldError{
				Field:   field,
				Code:    fieldErr.Tag(),
				Message: field + " failed validation rule " + fieldErr.Tag(),
			})
		}
		return errs
	}
	return nil
}
