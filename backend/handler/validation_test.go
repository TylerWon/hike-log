package handler

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

type validationTestSuite struct {
	suite.Suite
	validator *validator.Validate
}

func (suite *validationTestSuite) SetupTest() {
	v, err := registerCustomValidators()
	if err != nil {
		suite.T().Fatal("Error while registering custom validators: ", err)
	}

	suite.validator = v
}

func (suite *validationTestSuite) TestDivisibleByHalfValidator_ReturnsErrorForNonFloatField() {
	type sample struct {
		Value string `binding:"divisibleByHalf"`
	}

	err := suite.validator.Struct(sample{"test"})
	suite.NotNil(err)
}

func (suite *validationTestSuite) TestDivisibleByHalfValidator_ReturnsErrorForValueUndivisibleByPointFive() {
	type sample struct {
		Value float32 `binding:"divisibleByHalf"`
	}

	err := suite.validator.Struct(sample{2.3})
	suite.NotNil(err)
}

func (suite *validationTestSuite) TestDivisibleByHalfValidator_ValidationSucceeds() {
	type sample struct {
		Value float32 `binding:"divisibleByHalf"`
	}

	err := suite.validator.Struct(sample{0})
	suite.Nil(err)

	err = suite.validator.Struct(sample{0.5})
	suite.Nil(err)

	err = suite.validator.Struct(sample{1.0})
	suite.Nil(err)
}

func (suite *validationTestSuite) TestValidImageTypeValidator_ReturnsErrorForNonStringField() {
	type sample struct {
		Value int `binding:"validImageType"`
	}

	err := suite.validator.Struct(sample{2})
	suite.NotNil(err)
}

func (suite *validationTestSuite) TestValidImageTypeValidator_ReturnsErrorForNonImageType() {
	type sample struct {
		Value string `binding:"validImageType"`
	}

	err := suite.validator.Struct(sample{"text/html"})
	suite.NotNil(err)

	err = suite.validator.Struct(sample{"application/json"})
	suite.NotNil(err)

	err = suite.validator.Struct(sample{"video/mp4"})
	suite.NotNil(err)
}

func (suite *validationTestSuite) TestValidImageTypeValidator_ValidationSucceeds() {
	type sample struct {
		Value string `binding:"validImageType"`
	}

	err := suite.validator.Struct(sample{"image/jpeg"})
	suite.Nil(err)

	err = suite.validator.Struct(sample{"image/png"})
	suite.Nil(err)

	err = suite.validator.Struct(sample{"image/webp"})
	suite.Nil(err)

	err = suite.validator.Struct(sample{"image/heic"})
	suite.Nil(err)

	err = suite.validator.Struct(sample{"image/heif"})
	suite.Nil(err)
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, new(validationTestSuite))
}
