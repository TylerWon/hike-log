package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type dateTestSuite struct {
	suite.Suite
}

func (suite *dateTestSuite) TestMarshalJSON() {
	date := Date("2026-09-02")
	b, err := json.Marshal(date)
	suite.NoError(err)
	suite.Equal(`"2026-09-02"`, string(b))
}

func (suite *dateTestSuite) TestUnmarshalJSON() {
	var date Date
	err := json.Unmarshal([]byte(`"2026-09-02"`), &date)
	suite.NoError(err)
	suite.Equal(Date("2026-09-02"), date)
}

func TestDateTestSuite(t *testing.T) {
	suite.Run(t, new(dateTestSuite))
}
