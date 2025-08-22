package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/logger"
	"github.com/barretot/ifkpass/internal/util"
)

type IdentityProvider struct {
	client *cognitoidentityprovider.Client
}

var cfg = config.LoadConfig()

func NewIdentityProvider() IdentityProviderAdapter {
	awsCfg, _ := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(cfg.Region),
	)

	return &IdentityProvider{
		client: cognitoidentityprovider.NewFromConfig(awsCfg),
	}
}

func (idp *IdentityProvider) SignUp(ctx context.Context, email, password string) (string, error) {
	logger.Log.Info("starting cognito signup",
		"email", email,
	)

	hash := util.GenerateSecretHash(cfg.CognitoClientSecret, email, cfg.CognitoClientID)

	_, err := idp.client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(cfg.CognitoClientID),
		Username:   aws.String(email),
		Password:   aws.String(password),
		SecretHash: aws.String(hash),
	})

	if err != nil {
		logger.Log.Error("failed to sign up user in cognito", "email", email, "err", err)
		return "", fmt.Errorf("cognito signup: %w", err)
	}

	response, err := idp.client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		Username:   aws.String(email),
		UserPoolId: aws.String(cfg.CognitoUserPoolID),
	})

	if err != nil {
		logger.Log.Error("failed to get user id from cognito", "email", email, "err", err)
		return "", fmt.Errorf("get user id from cognito: %w", err)
	}

	logger.Log.Info("user signed up successfully in cognito", "email", email)
	logger.Log.Info("retrieved user id from cognito", "email", email, "user_id", *response.Username)

	return *response.Username, nil

}

func (idp *IdentityProvider) SignIn(ctx context.Context, email, password string) (*string, error) {
	logger.Log.Info("starting cognito signin", "email", email)

	hash := util.GenerateSecretHash(cfg.CognitoClientSecret, email, cfg.CognitoClientID)

	r, err := idp.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(cfg.CognitoClientID),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": hash,
		},
	})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			logger.Log.Warn("cognito signin api error", "code", ae.ErrorCode(), "email", email)

			switch ae.ErrorCode() {
			case "NotAuthorizedException":
				return nil, fmt.Errorf("not authorized: %s", ae.ErrorMessage())
			case "UserNotFoundException":
				return nil, fmt.Errorf("user not found: %s", ae.ErrorMessage())
			case "UserNotConfirmedException":
				return nil, fmt.Errorf("user not confirmed: %s", ae.ErrorMessage())
			default:
				return nil, fmt.Errorf("cognito error %s: %s", ae.ErrorCode(), ae.ErrorMessage())
			}
		}
		return nil, fmt.Errorf("cognito initiate auth failed: %w", err)
	}

	if r.ChallengeName != "" {
		switch r.ChallengeName {
		case types.ChallengeNameTypeSmsMfa, types.ChallengeNameTypeSoftwareTokenMfa:
			return nil, fmt.Errorf("mfa required: %s", r.ChallengeName)
		case types.ChallengeNameTypeNewPasswordRequired:
			return nil, fmt.Errorf("new password required")
		default:
			return nil, fmt.Errorf("unsupported cognito challenge: %s", r.ChallengeName)
		}
	}

	if r.AuthenticationResult == nil || r.AuthenticationResult.AccessToken == nil {
		logger.Log.Error("missing authentication result", "email", email)
		return nil, fmt.Errorf("cognito signin: missing authentication result")
	}

	logger.Log.Info("user signed in successfully in cognito", "email", email)
	return r.AuthenticationResult.AccessToken, nil
}

func (identityprovider *IdentityProvider) GetUserId(ctx context.Context, email string) (string, error) {
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

func (identityprovider *IdentityProvider) IsEmailVerified(ctx context.Context, email string) (bool, error) {
	response, err := identityprovider.client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		Username:   aws.String(email),
		UserPoolId: aws.String(cfg.CognitoUserPoolID),
	})

	if err != nil {
		logger.Log.Error("failed to get user from cognito", "email", email, "err", err)
		return false, fmt.Errorf("get user from cognito: %w", err)
	}

	var emailVerified bool
	for _, attr := range response.UserAttributes {
		if *attr.Name == "email_verified" {
			emailVerified = *attr.Value == "true"
			break
		}
	}

	logger.Log.Info("checked email verification status", "email", email, "verified", emailVerified)
	return emailVerified, nil
}

func (identityprovider *IdentityProvider) ConfirmEmail(ctx context.Context, email, code string) error {
	hash := util.GenerateSecretHash(cfg.CognitoClientSecret, email, cfg.CognitoClientID)

	_, err := identityprovider.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(cfg.CognitoClientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
		SecretHash:       aws.String(hash),
	})

	if err != nil {
		logger.Log.Error("failed to confirm email", "email", email, "err", err)
		return fmt.Errorf("cognito signup: %w", err)
	}

	logger.Log.Info("user confirm email sucessfully", "email", email)
	return nil
}
