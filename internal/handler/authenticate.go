package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/util"
)

func HandleAuthenticate(ctx context.Context, event events.APIGatewayProxyRequest, cfg config.AppConfig) (events.APIGatewayProxyResponse, error) {
	return util.NewSuccessResponse(404, "TODO IMPLEMENTED"), nil
}
