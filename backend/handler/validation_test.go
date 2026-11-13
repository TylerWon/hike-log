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
	v, err := RegisterCustomValidators()
	if err != nil {
		suite.T().Fatal("Error while registering custom validators: ", err)
	}

	suite.validator = v
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, new(validationTestSuite))
}

func (suite *validationTestSuite) TestDivisibleByHalfValidator_ReturnsErrorForNonFloatField() {
	type sample struct {
		Value string `binding:"divisibleByHalf"`
	}

	err := suite.validator.Struct(sample{"test"})
	suite.Require().Error(err)
}

func (suite *validationTestSuite) TestDivisibleByHalfValidator_ReturnsErrorForValueUndivisibleByHalf() {
	type sample struct {
		Value float32 `binding:"divisibleByHalf"`
	}

	err := suite.validator.Struct(sample{2.3})
	suite.Require().Error(err)
}

func (suite *validationTestSuite) TestDivisibleByHalfValidator_ValidationSucceeds() {
	type sample struct {
		Value float32 `binding:"divisibleByHalf"`
	}

	err := suite.validator.Struct(sample{0})
	suite.Require().NoError(err)

	err = suite.validator.Struct(sample{0.5})
	suite.Require().NoError(err)

	err = suite.validator.Struct(sample{1.0})
	suite.Require().NoError(err)
}

func (suite *validationTestSuite) TestValidImageTypeValidator_ReturnsErrorForNonStringField() {
	type sample struct {
		Value int `binding:"validImageType"`
	}

	err := suite.validator.Struct(sample{2})
	suite.Require().Error(err)
}

func (suite *validationTestSuite) TestValidImageTypeValidator_ReturnsErrorForInvalidImageType() {
	type sample struct {
		Value string `binding:"validImageType"`
	}

	err := suite.validator.Struct(sample{"text/html"})
	suite.Require().Error(err)

	err = suite.validator.Struct(sample{"application/json"})
	suite.Require().Error(err)

	err = suite.validator.Struct(sample{"video/mp4"})
	suite.Require().Error(err)
}

func (suite *validationTestSuite) TestValidImageTypeValidator_ValidationSucceeds() {
	type sample struct {
		Value string `binding:"validImageType"`
	}

	err := suite.validator.Struct(sample{"image/jpeg"})
	suite.Require().NoError(err)

	err = suite.validator.Struct(sample{"image/png"})
	suite.Require().NoError(err)

	err = suite.validator.Struct(sample{"image/webp"})
	suite.Require().NoError(err)

	err = suite.validator.Struct(sample{"image/heic"})
	suite.Require().NoError(err)

	err = suite.validator.Struct(sample{"image/heif"})
	suite.Require().NoError(err)
}

func (suite *validationTestSuite) TestValidObjectKeyValidator_ReturnsErrorForNonStringField() {
	type sample struct {
		Value bool `binding:"validObjectKey"`
	}

	err := suite.validator.Struct(sample{true})
	suite.Require().Error(err)
}

func (suite *validationTestSuite) TestValidObjectKeyValidator_ReturnsErrorForInvalidObjectKey() {
	type sample struct {
		Value string `binding:"validObjectKey"`
	}

	err := suite.validator.Struct(sample{"hikes/"})
	suite.Require().Error(err)

	err = suite.validator.Struct(sample{"hikes/asgbsdfg123/"})
	suite.Require().Error(err)

	err = suite.validator.Struct(sample{"hikes/1/files/"})
	suite.Require().Error(err)

	err = suite.validator.Struct(sample{"hikes/1/photos/1"})
	suite.Require().Error(err)

	err = suite.validator.Struct(sample{"hikes/1/photos/asdf091123asdf"})
	suite.Require().Error(err)
}

func (suite *validationTestSuite) TestValidObjectKeyValidator_ValidationSucceeds() {
	type sample struct {
		Value string `binding:"validObjectKey"`
	}

	err := suite.validator.Struct(sample{"hikes/1/photos/acde070d-8c4c-4f0d-9d8a-162843c10333"})
	suite.Require().NoError(err)
}
