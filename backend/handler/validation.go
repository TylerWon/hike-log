package handler

import (
	"errors"
	"math"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var VALID_IMAGE_TYPES = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/heic": {},
	"image/heif": {},
}

// Registers custom struct field validators used by Gin for model validation
func registerCustomValidators() (*validator.Validate, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("Failed to register custom validators")
	}

	v.RegisterValidation("divisibleByHalf", divisibleByHalfValidator)
	v.RegisterValidation("validImageType", validImageTypeValidator)

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

// Checks that a string field is a valid image MIME type
func validImageTypeValidator(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}
	_, ok := VALID_IMAGE_TYPES[fl.Field().String()]
	return ok
}
