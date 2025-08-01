package config

import "os"

type AppConfig struct {
	Region    string
	TableName string
}

func LoadConfig() AppConfig {
	return AppConfig{
		Region:    os.Getenv("AWS_REGION"),
		TableName: os.Getenv("DYNAMO_TABLE_NAME"),
	}
}
