package handler

import (
	"errors"
	"math"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Registers custom struct field validators to be used by Gin for model validation
func registerCustomValidators() (*validator.Validate, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("Failed to register custom validators")
	}

	v.RegisterValidation("divisibleByHalf", divisibleByHalfValidator)

	return v, nil
}

// Checks that a float field is divisible by 0.5
func divisibleByHalfValidator(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Float32, reflect.Float64:
		doubled := fl.Field().Float() * 2
		_, frac := math.Modf(doubled)
		return frac < 1e-9 || frac > 1-1e-9
	default:
		return false
	}
}
