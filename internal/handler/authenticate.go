package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/barretot/ifkpass/internal/apperrors"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/contextkeys"
	"github.com/barretot/ifkpass/internal/dto"
	"github.com/barretot/ifkpass/internal/identity"
	"github.com/barretot/ifkpass/internal/logger"
	"github.com/barretot/ifkpass/internal/service"
	"github.com/barretot/ifkpass/internal/store/dynamostore"
	"github.com/barretot/ifkpass/internal/util"
	"github.com/barretot/ifkpass/internal/validator"
)

func HandleAuthenticate(ctx context.Context, event events.APIGatewayProxyRequest, cfg config.AppConfig) (events.APIGatewayProxyResponse, error) {
	var input dto.AuthenticateInput

	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		return util.NewErrorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	if err := validator.ValidateRequest(input); err != nil {
		return util.NewErrorResponse(http.StatusBadRequest, "validation error"), nil
	}

	repo := dynamostore.NewDynamoProfileRepository(cfg)
	identityprovider := identity.NewIdentityProvider(cfg)
	s := service.NewAuthenticateService(repo, identityprovider)

	requestID, _ := ctx.Value(contextkeys.RequestID).(string)

	logger.Log.Info("received create user request",
		"email", input.Email,
		"request_id", requestID,
	)

	token, err := s.Authenticate(ctx, cfg, input.Email, input.Password)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return util.EncodeJson(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		}

		return util.EncodeJson(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return util.EncodeJson(http.StatusOK, map[string]any{
		"token": token,
	})
}
