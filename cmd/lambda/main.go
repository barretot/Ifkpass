package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/contextkeys"
	"github.com/barretot/ifkpass/internal/handler"
	"github.com/barretot/ifkpass/internal/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.Init(cfg.GoEnv)

	lambda.Start(func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		requestID := event.RequestContext.RequestID

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, contextkeys.RequestID, requestID)

		switch event.Resource {
		case "/user":
			if event.HTTPMethod == "POST" {
				return handler.HandleCreateUser(ctx, event, cfg)
			}
		case "/auth":
			if event.HTTPMethod == "POST" {
				return handler.HandleAuthenticate(ctx, event, cfg)
			}
		case "/verify":
			if event.HTTPMethod == "POST" {
				return handler.HandleVerifyEmail(ctx, event, cfg)
			}
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"message":"Route not found"}`,
		}, nil
	})
}
