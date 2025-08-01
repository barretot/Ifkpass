package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/dto"
	"github.com/barretot/ifkpass/internal/repo"
	"github.com/barretot/ifkpass/internal/service"
	"github.com/barretot/ifkpass/internal/util"
	"github.com/barretot/ifkpass/internal/validator"
)

func HandleCreateUser(ctx context.Context, event events.APIGatewayProxyRequest, cfg config.AppConfig) (events.APIGatewayProxyResponse, error) {
	var input dto.CreateUserInput

	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		return util.NewErrorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	if err := validator.ValidateRequest(input); err != nil {
		return util.NewErrorResponse(http.StatusBadRequest, "validation error"), nil
	}

	repo := repo.NewDynamoUserRepository(cfg)

	userService := service.NewUserService(repo)

	if err := userService.CreateUser(ctx, input.Name, input.Email); err != nil {
		return util.NewErrorResponse(500, err.Error()), nil
	}

	return util.NewSuccessResponse(201, "user created"), nil
}
