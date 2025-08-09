package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Region    string
	TableName string
}

func LoadConfig() AppConfig {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return AppConfig{
		Region:    os.Getenv("AWS_REGION"),
		TableName: os.Getenv("PROFILES_TABLE_NAME"),
	}
}
