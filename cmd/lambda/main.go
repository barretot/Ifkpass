package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/handler"
)

func main() {
	appConfig := config.LoadConfig()

	lambda.Start(func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch event.Resource {
		case "/user":
			if event.HTTPMethod == "POST" {
				return handler.HandleCreateUser(ctx, event, appConfig)
			}
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"message":"Route not found"}`,
		}, nil
	})
}
