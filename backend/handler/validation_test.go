package handler

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

type validationTestSuite struct {
	suite.Suite
	v *validator.Validate
}

func (suite *validationTestSuite) SetupTest() {
	v, err := registerCustomValidators()
	if err != nil {
		suite.T().Fatal("Error while registering custom validators: ", err)
	}

	suite.v = v
}

func (suite *validationTestSuite) TestHalfStepValidator_ReturnsErrorForNonFloatField() {
	type sample struct {
		Value string `binding:"halfstep"`
	}

	err := suite.v.Struct(sample{"test"})
	suite.NotNil(err)
}

func (suite *validationTestSuite) TestHalfStepValidator_ReturnsErrorForValueUndivisibleByPointFive() {
	type sample struct {
		Value float32 `binding:"halfstep"`
	}

	err := suite.v.Struct(sample{2.3})
	suite.NotNil(err)
}

func (suite *validationTestSuite) TestHalfStepValidator_ValidationSucceeds() {
	type sample struct {
		Value float32 `binding:"halfstep"`
	}

	err := suite.v.Struct(sample{0})
	suite.Nil(err)

	err = suite.v.Struct(sample{0.5})
	suite.Nil(err)

	err = suite.v.Struct(sample{1.0})
	suite.Nil(err)
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, new(validationTestSuite))
}
