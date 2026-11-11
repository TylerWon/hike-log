package config_test

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/stretchr/testify/suite"
)

type configTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}

func setEnvVars(t *testing.T, env string) {
	t.Setenv("ENV", env)
	t.Setenv("FRONTEND_URL", "http://localhost:5173")
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "hike_log_test_db")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "password")
	t.Setenv("AWS_S3_BUCKET_NAME", "hike-log-test-bucket")
	t.Setenv("AWS_ACCESS_KEY_ID", "localstack")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("SESSION_SECRET", "vDuTmQE8HJj_SsRT32y8Wz8yH97FoRCT5VUiNJlfsGI")
	t.Setenv("ADMIN_PASSWORD_HASH", "$2a$10$haTJ.11uqjdIRUD4yxMSt.HfJiOL9.WjyjH1D2C.qv8okP0c0wuXC")

	if env == "local" {
		t.Setenv("AWS_ENDPOINT_URL", "http://localstack:4566")
		t.Setenv("AWS_EXTERNAL_ENDPOINT_URL", "http://localhost:4566")
	}
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenEnvIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("ENV", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenFrontendURLIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("FRONTEND_URL", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenDBHostIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("DB_HOST", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenDBPortIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("DB_PORT", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenDBNameIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("DB_NAME", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenDBUserIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("DB_USER", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenDBPasswordIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("DB_PASSWORD", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenAWSS3BucketNameIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("AWS_S3_BUCKET_NAME", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenAWSAccessKeyIDIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("AWS_ACCESS_KEY_ID", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenAWSSecretAccessKeyIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("AWS_SECRET_ACCESS_KEY", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenAWSRegionIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("AWS_REGION", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenSessionSecretIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("SESSION_SECRET", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenAdminPasswordHashIsMissing() {
	setEnvVars(suite.T(), "production")
	suite.T().Setenv("ADMIN_PASSWORD_HASH", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenAWSEndpointURLIsMissingInLocalEnv() {
	setEnvVars(suite.T(), "local")
	suite.T().Setenv("AWS_ENDPOINT_URL", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsErrorWhenAWSExternalEndpointURLIsMissingInLocalEnv() {
	setEnvVars(suite.T(), "local")
	suite.T().Setenv("AWS_EXTERNAL_ENDPOINT_URL", "")

	cfg, err := config.New()

	suite.Require().Nil(cfg)
	suite.Require().Error(err)
}

func (suite *configTestSuite) TestNew_ReturnsConfigForNonLocalEnv() {
	setEnvVars(suite.T(), "production")

	cfg, err := config.New()

	expected := &config.Config{
		Env:         "production",
		FrontendURL: "http://localhost:5173",
		Database: config.DBConfig{
			DBHost:     "db",
			DBPort:     "5432",
			DBName:     "hike_log_test_db",
			DBUser:     "postgres",
			DBPassword: "password",
		},
		AWS: config.AWSConfig{
			AccessKeyID:         "localstack",
			SecretAccessKey:     "localstack",
			Region:              "us-east-1",
			EndpointURL:         "",
			ExternalEndpointURL: "",
			S3: config.S3Config{
				BucketName: "hike-log-test-bucket",
			},
		},
		Auth: config.AuthConfig{
			SessionSecret:     "vDuTmQE8HJj_SsRT32y8Wz8yH97FoRCT5VUiNJlfsGI",
			AdminPasswordHash: "$2a$10$haTJ.11uqjdIRUD4yxMSt.HfJiOL9.WjyjH1D2C.qv8okP0c0wuXC",
		},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expected, cfg)
}

func (suite *configTestSuite) TestNew_ReturnsConfigForLocalEnv() {
	setEnvVars(suite.T(), "local")

	cfg, err := config.New()

	expected := &config.Config{
		Env:         "local",
		FrontendURL: "http://localhost:5173",
		Database: config.DBConfig{
			DBHost:     "db",
			DBPort:     "5432",
			DBName:     "hike_log_test_db",
			DBUser:     "postgres",
			DBPassword: "password",
		},
		AWS: config.AWSConfig{
			AccessKeyID:         "localstack",
			SecretAccessKey:     "localstack",
			Region:              "us-east-1",
			EndpointURL:         "http://localstack:4566",
			ExternalEndpointURL: "http://localhost:4566",
			S3: config.S3Config{
				BucketName: "hike-log-test-bucket",
			},
		},
		Auth: config.AuthConfig{
			SessionSecret:     "vDuTmQE8HJj_SsRT32y8Wz8yH97FoRCT5VUiNJlfsGI",
			AdminPasswordHash: "$2a$10$haTJ.11uqjdIRUD4yxMSt.HfJiOL9.WjyjH1D2C.qv8okP0c0wuXC",
		},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expected, cfg)
}
