package util

import (
    "encoding/json"

    "github.com/aws/aws-lambda-go/events"
)

func NewSuccessResponse(code int, message string) events.APIGatewayProxyResponse {
    body, _ := json.Marshal(map[string]string{"message": message})
    return events.APIGatewayProxyResponse{
        StatusCode: code,
        Body:       string(body),
    }
}

func NewErrorResponse(code int, message string) events.APIGatewayProxyResponse {
    body, _ := json.Marshal(map[string]string{"error": message})
    return events.APIGatewayProxyResponse{
        StatusCode: code,
        Body:       string(body),
    }
}