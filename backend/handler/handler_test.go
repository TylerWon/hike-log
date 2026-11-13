package handler_test

import (
	"errors"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (suite *handlerTestSuite) TestNewHandler_ReturnsErrorWhenRegisterValidatorsFails() {
	mockValidatorManager := &handler.MockValidatorManager{
		RegisterValidatorsError: errors.New("Something went wrong"),
	}

	h, err := handler.NewHandler(
		&store.MockStore{},
		&photostorage.MockPhotoStorage{},
		&auth.MockAuth{},
		mockValidatorManager,
	)

	suite.Require().Error(err)
	suite.Require().Nil(h)
}

func (suite *handlerTestSuite) TestNewHandler_CreatesHandler() {
	h, err := handler.NewHandler(
		&store.MockStore{},
		&photostorage.MockPhotoStorage{},
		&auth.MockAuth{},
		&handler.MockValidatorManager{},
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(h)
	suite.Require().NotNil(h.Auth)
	suite.Require().NotNil(h.Hikes)
	suite.Require().NotNil(h.Photos)
}
