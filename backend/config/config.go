package config

import (
	"fmt"
	"os"
)

// Config stores configuration for the app.
type Config struct {
	Database DBConfig
	AWS      AWSConfig
	Auth     AuthConfig
	// The current deployment environment (e.g. local, staging, production).
	Env string
	// The URL of the frontend app.
	FrontendURL string
}

// AuthConfig stores all auth related configuration.
type AuthConfig struct {
	// The secret used to sign session cookies.
	SessionSecret string
	// A bcrypt hash of the admin password.
	AdminPasswordHash string
}

// AWSConfig stores all AWS related configuration.
type AWSConfig struct {
	// The AWS access key ID.
	AccessKeyID string
	// The AWS secret access key.
	SecretAccessKey string
	// The AWS region.
	Region string
	// (In local environment only) The AWS API endpoint accessible from outside Docker.
	ExternalEndpointURL string
	// (In local environment only) The AWS API endpoint.
	EndpointURL string
	// S3 configuration
	S3 S3Config
}

// S3Config stores all AWS S3 related configuration.
type S3Config struct {
	// The name of the S3 bucket used to store hike data.
	BucketName string
}

// DBConfig stores all database related configuration.
type DBConfig struct {
	// The hostname of the Postgres server.
	DBHost string
	// The port of the Postgres server.
	DBPort string
	// The name of the Postgres database.
	DBName string
	// The username to access the database.
	DBUser string
	// The password to access the database.
	DBPassword string
}

// Constructs a Config from environment variables.
func New() (*Config, error) {
	// Values required in all environments
	env, err := getEnvVar("ENV")
	if err != nil {
		return nil, err
	}
	frontendURL, err := getEnvVar("FRONTEND_URL")
	if err != nil {
		return nil, err
	}
	dbHost, err := getEnvVar("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := getEnvVar("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnvVar("DB_NAME")
	if err != nil {
		return nil, err
	}
	dbUser, err := getEnvVar("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPassword, err := getEnvVar("DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	bucketName, err := getEnvVar("AWS_S3_BUCKET_NAME")
	if err != nil {
		return nil, err
	}
	accessKeyID, err := getEnvVar("AWS_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}
	secretAccessKey, err := getEnvVar("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}
	region, err := getEnvVar("AWS_REGION")
	if err != nil {
		return nil, err
	}
	sessionSecret, err := getEnvVar("SESSION_SECRET")
	if err != nil {
		return nil, err
	}
	adminPasswordHash, err := getEnvVar("ADMIN_PASSWORD_HASH")
	if err != nil {
		return nil, err
	}

	// Values required in local environment only
	var endpointURL, externalEndpointURL string
	if env == "local" {
		endpointURL, err = getEnvVar("AWS_ENDPOINT_URL")
		if err != nil {
			return nil, err
		}
		externalEndpointURL, err = getEnvVar("AWS_EXTERNAL_ENDPOINT_URL")
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		Env:         env,
		FrontendURL: frontendURL,
		Database: DBConfig{
			DBHost:     dbHost,
			DBPort:     dbPort,
			DBName:     dbName,
			DBUser:     dbUser,
			DBPassword: dbPassword,
		},
		AWS: AWSConfig{
			AccessKeyID:         accessKeyID,
			SecretAccessKey:     secretAccessKey,
			Region:              region,
			EndpointURL:         endpointURL,
			ExternalEndpointURL: externalEndpointURL,
			S3: S3Config{
				BucketName: bucketName,
			},
		},
		Auth: AuthConfig{
			SessionSecret:     sessionSecret,
			AdminPasswordHash: adminPasswordHash,
		},
	}, nil
}

func getEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s environment variable is missing or empty", key)
	}
	return value, nil
}
