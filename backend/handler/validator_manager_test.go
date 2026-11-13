package handler_test

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

type validatorManagerTestSuite struct {
	suite.Suite
	validate *validator.Validate
}

func (suite *validatorManagerTestSuite) SetupTest() {
	validatorManager, err := handler.NewValidatorManager()
	if err != nil {
		suite.T().Fatal("Error while creating validator: ", err)
	}

	err = validatorManager.RegisterValidators()
	if err != nil {
		suite.T().Fatal("Error while registering validators: ", err)
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		suite.T().Fatal("Failed to obtain validator engine")
	}

	suite.validate = v
}

func TestValidatorManagerTestSuite(t *testing.T) {
	suite.Run(t, new(validatorManagerTestSuite))
}

func (suite *validatorManagerTestSuite) TestNewValidatorManager_CreatesValidatorManager() {
	validatorManager, err := handler.NewValidatorManager()
	suite.Require().NoError(err)
	suite.Require().NotNil(validatorManager)
}

func (suite *validatorManagerTestSuite) TestDivisibleByHalfValidator_ReturnsErrorForNonFloatField() {
	type sample struct {
		Value string `binding:"divisibleByHalf"`
	}

	err := suite.validate.Struct(sample{"test"})
	suite.Require().Error(err)
}

func (suite *validatorManagerTestSuite) TestDivisibleByHalfValidator_ReturnsErrorForValueUndivisibleByHalf() {
	type sample struct {
		Value float32 `binding:"divisibleByHalf"`
	}

	err := suite.validate.Struct(sample{2.3})
	suite.Require().Error(err)
}

func (suite *validatorManagerTestSuite) TestDivisibleByHalfValidator_ValidationSucceeds() {
	type sample struct {
		Value float32 `binding:"divisibleByHalf"`
	}

	err := suite.validate.Struct(sample{0})
	suite.Require().NoError(err)

	err = suite.validate.Struct(sample{0.5})
	suite.Require().NoError(err)

	err = suite.validate.Struct(sample{1.0})
	suite.Require().NoError(err)
}

func (suite *validatorManagerTestSuite) TestValidImageTypeValidator_ReturnsErrorForNonStringField() {
	type sample struct {
		Value int `binding:"validImageType"`
	}

	err := suite.validate.Struct(sample{2})
	suite.Require().Error(err)
}

func (suite *validatorManagerTestSuite) TestValidImageTypeValidator_ReturnsErrorForInvalidImageType() {
	type sample struct {
		Value string `binding:"validImageType"`
	}

	err := suite.validate.Struct(sample{"text/html"})
	suite.Require().Error(err)

	err = suite.validate.Struct(sample{"application/json"})
	suite.Require().Error(err)

	err = suite.validate.Struct(sample{"video/mp4"})
	suite.Require().Error(err)
}

func (suite *validatorManagerTestSuite) TestValidImageTypeValidator_ValidationSucceeds() {
	type sample struct {
		Value string `binding:"validImageType"`
	}

	err := suite.validate.Struct(sample{"image/jpeg"})
	suite.Require().NoError(err)

	err = suite.validate.Struct(sample{"image/png"})
	suite.Require().NoError(err)

	err = suite.validate.Struct(sample{"image/webp"})
	suite.Require().NoError(err)

	err = suite.validate.Struct(sample{"image/heic"})
	suite.Require().NoError(err)

	err = suite.validate.Struct(sample{"image/heif"})
	suite.Require().NoError(err)
}

func (suite *validatorManagerTestSuite) TestValidObjectKeyValidator_ReturnsErrorForNonStringField() {
	type sample struct {
		Value bool `binding:"validObjectKey"`
	}

	err := suite.validate.Struct(sample{true})
	suite.Require().Error(err)
}

func (suite *validatorManagerTestSuite) TestValidObjectKeyValidator_ReturnsErrorForInvalidObjectKey() {
	type sample struct {
		Value string `binding:"validObjectKey"`
	}

	err := suite.validate.Struct(sample{"hikes/"})
	suite.Require().Error(err)

	err = suite.validate.Struct(sample{"hikes/asgbsdfg123/"})
	suite.Require().Error(err)

	err = suite.validate.Struct(sample{"hikes/1/files/"})
	suite.Require().Error(err)

	err = suite.validate.Struct(sample{"hikes/1/photos/1"})
	suite.Require().Error(err)

	err = suite.validate.Struct(sample{"hikes/1/photos/asdf091123asdf"})
	suite.Require().Error(err)
}

func (suite *validatorManagerTestSuite) TestValidObjectKeyValidator_ValidationSucceeds() {
	type sample struct {
		Value string `binding:"validObjectKey"`
	}

	err := suite.validate.Struct(sample{"hikes/1/photos/acde070d-8c4c-4f0d-9d8a-162843c10333"})
	suite.Require().NoError(err)
}
