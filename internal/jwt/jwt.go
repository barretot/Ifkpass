package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/barretot/ifkpass/internal/config"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var cfg = config.LoadConfig()

var (
	baseUrl    = cfg.CognitoURL
	userPoolId = cfg.CognitoUserPoolID
	clientId   = cfg.CognitoClientID
)

var jwksUrl = fmt.Sprintf("%s/%s/.well-known/jwks.json", baseUrl, userPoolId)

func VerifyToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("unauthorized")
	}

	rawToken := strings.TrimPrefix(authHeader, "Bearer ")

	set, err := jwk.Fetch(context.Background(), jwksUrl)
	if err != nil {
		return "", fmt.Errorf("unauthorized: %w", err)
	}

	token, err := jwt.ParseString(
		rawToken, jwt.WithKeySet(set),
		jwt.WithValidate(true),
		// jwt.WithAudience(clientId),
		jwt.WithIssuer(fmt.Sprintf("%s/%s", baseUrl, userPoolId)),
	)

	if err != nil {
		fmt.Printf("JWT error: %v\n", err)
		return "", fmt.Errorf("token inválido ou expirado: %w", err)
	}
	subject := token.Subject()

	return subject, nil
}
