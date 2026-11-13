package testutils

import "github.com/TylerWon/hike-log/backend/config"

var testConfig = config.Config{
	Env:         "local",
	FrontendURL: "http://localhost:5173",
	// Configuration for Docker database container
	Database: config.DBConfig{
		DBHost:     "db",
		DBPort:     "5432",
		DBName:     "hike_log_test_db",
		DBUser:     "postgres",
		DBPassword: "password",
	},
	// Configuration for Docker LocalStack container
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
		SessionSecret: "vDuTmQE8HJj_SsRT32y8Wz8yH97FoRCT5VUiNJlfsGI",
		// bcrypt hash for "admin"
		AdminPasswordHash: "$2a$10$haTJ.11uqjdIRUD4yxMSt.HfJiOL9.WjyjH1D2C.qv8okP0c0wuXC",
	},
}

// Returns a copy of the default Config to use for testing.
func NewTestConfig() config.Config {
	return testConfig
}
