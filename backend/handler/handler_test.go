package handler_test

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/database"
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

func (suite *handlerTestSuite) TestNewHandler_CreatesHandler() {
	h, err := handler.NewHandler(
		&database.MockDatabase{},
		&store.MockStore{},
		&photostorage.MockPhotoStorage{},
		&auth.MockAuth{},
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(h)
	suite.Require().NotNil(h.Auth)
	suite.Require().NotNil(h.Health)
	suite.Require().NotNil(h.Hikes)
	suite.Require().NotNil(h.Photos)
}
