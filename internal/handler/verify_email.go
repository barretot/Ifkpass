package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/dto"
	"github.com/barretot/ifkpass/internal/util"
	"github.com/barretot/ifkpass/internal/validator"
)

func HandleVerifyEmail(ctx context.Context, event events.APIGatewayProxyRequest, cfg config.AppConfig) (events.APIGatewayProxyResponse, error) {
	var input dto.VerifyEmailInput

	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		return util.NewErrorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	if err := validator.ValidateRequest(input); err != nil {
		return util.NewErrorResponse(http.StatusBadRequest, "validation error"), nil
	}

	return util.NewSuccessResponse(404, "TODO IMPLEMENTED"), nil
}
