package identity

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/logger"
	"github.com/barretot/ifkpass/internal/util"
)

type IdentityProvider struct {
	client *cognitoidentityprovider.Client
}

func NewIdentityProvider(cfg config.AppConfig) IdentityProviderAdapter {
	awsCfg, _ := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(cfg.Region),
	)

	return &IdentityProvider{
		client: cognitoidentityprovider.NewFromConfig(awsCfg),
	}
}

func (identityprovider *IdentityProvider) SignUp(ctx context.Context, cfg config.AppConfig, email, password string) error {
	logger.Log.Info("starting cognito signup",
		"email", email,
	)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	hash := util.GenerateSecretHash(cfg.CognitoClientSecret, email, cfg.CognitoClientID)

	_, err := identityprovider.client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(cfg.CognitoClientID),
		Username:   aws.String(email),
		Password:   aws.String(password),
		SecretHash: aws.String(hash),
	})
	if err != nil {
		logger.Log.Error("failed to sign up user in cognito", "email", email, "err", err)
		return fmt.Errorf("cognito signup: %w", err)
	}

	logger.Log.Info("user signed up successfully in cognito", "email", email)
	return nil
}

func (identityprovider *IdentityProvider) GetUserId(ctx context.Context, cfg config.AppConfig, email string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	response, err := identityprovider.client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		Username:   aws.String(email),
		UserPoolId: aws.String(cfg.CognitoUserPoolID),
	})
	if err != nil {
		logger.Log.Error("failed to get user id from cognito", "email", email, "err", err)
		return "", fmt.Errorf("get user id from cognito: %w", err)
	}

	logger.Log.Info("retrieved user id from cognito", "email", email, "user_id", *response.Username)
	return *response.Username, nil
}
