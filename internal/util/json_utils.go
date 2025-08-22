package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func EncodeJson[T any](statusCode int, data T) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // 👈 evita \u0026 no output

	if err := encoder.Encode(data); err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to encode json: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		Headers:         map[string]string{"Content-Type": "application/json; charset=utf-8"},
		Body:            strings.TrimSpace(buf.String()),
		IsBase64Encoded: false,
	}, nil
}

func DecodeJson[T any](req events.APIGatewayProxyRequest) (T, error) {
	var data T

	body := req.Body
	if req.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return data, fmt.Errorf("decode base64 body: %w", err)
		}
		body = string(decoded)
	}

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return data, fmt.Errorf("decode json failed: %w", err)
	}

	return data, nil
}
