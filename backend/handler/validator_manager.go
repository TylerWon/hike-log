package handler

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var VALID_IMAGE_TYPES = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/heic": {},
	"image/heif": {},
}

// ValidatorManager manages custom struct field validators. These are used by Gin for model validation.
type ValidatorManager interface {
	RegisterValidators() error
}

// validatorManagerImpl is an implementation of the ValidatorManager interface.
type validatorManagerImpl struct {
	validate *validator.Validate
}

// Creates a new ValidatorManager.
func NewValidatorManager() (ValidatorManager, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("Failed to obtain validator engine")
	}

	return &validatorManagerImpl{validate: v}, nil
}

// Registers custom struct field validators.
func (c *validatorManagerImpl) RegisterValidators() error {
	err := c.validate.RegisterValidation("divisibleByHalf", divisibleByHalfValidator)
	if err != nil {
		return fmt.Errorf("Failed to register divisibleByHalf validator: %w", err)
	}

	err = c.validate.RegisterValidation("validImageType", validImageTypeValidator)
	if err != nil {
		return fmt.Errorf("Failed to register validImageType validator: %w", err)
	}

	err = c.validate.RegisterValidation("validObjectKey", validObjectKeyValidator)
	if err != nil {
		return fmt.Errorf("Failed to register validObjectKey validator: %w", err)
	}

	return nil
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

// Checks that a string field is a valid S3 object key. Format is "hikes/<hike_id>/photos/<uuid>".
func validObjectKeyValidator(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}

	parts := strings.Split(fl.Field().String(), "/")
	if len(parts) != 4 || parts[0] != "hikes" {
		return false
	}

	_, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return false
	}

	if parts[2] != "photos" {
		return false
	}

	_, err = uuid.Parse(parts[3])
	if err != nil {
		return false
	}

	return true
}
