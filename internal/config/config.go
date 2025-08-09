package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Region           string
	GoEnv            string
	BcryptSaltRounds string

	UsersTableName    string
	ProfilesTableName string

	ResendMailAPIKey string

	CognitoURL          string
	CognitoClientID     string
	CognitoClientSecret string
	CognitoUserPoolID   string

	ProfileBucketName string
}

func LoadConfig() AppConfig {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return AppConfig{
		ProfilesTableName: os.Getenv("PROFILES_TABLE_NAME"),
		Region:            os.Getenv("AWS_REGION"),
		GoEnv:             os.Getenv("GO_ENV"),

		ResendMailAPIKey: os.Getenv("RESEND_MAIL_API_KEY"),

		CognitoURL:          os.Getenv("COGNITO_URL"),
		CognitoClientID:     os.Getenv("COGNITO_CLIENT_ID"),
		CognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		CognitoUserPoolID:   os.Getenv("COGNITO_USER_POOL_ID"),

		ProfileBucketName: os.Getenv("PROFILE_BUCKET_NAME"),
	}
}
