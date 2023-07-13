package utils

import (
	"fmt"

	"github.com/go-playground/validator"
)

type ErrorResponse struct {
    FailedField string
    Tag         string
    Value       string
}

func ValidateJSON[T any](b T) ([]*ErrorResponse, error) {
	var errors []*ErrorResponse
	validate := validator.New()
	if err := validate.Struct(&b); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
            var element ErrorResponse
            element.FailedField = err.StructNamespace()
            element.Tag = err.Tag()
            element.Value = err.Param()
            errors = append(errors, &element)
        }
		return errors, fmt.Errorf("validation failed")
	}
	return nil, nil
}

