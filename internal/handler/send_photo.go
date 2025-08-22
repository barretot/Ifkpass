package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/contextkeys"
	"github.com/barretot/ifkpass/internal/dto"
	"github.com/barretot/ifkpass/internal/jwt"
	"github.com/barretot/ifkpass/internal/logger"
	"github.com/barretot/ifkpass/internal/service"
	"github.com/barretot/ifkpass/internal/storage"
	"github.com/barretot/ifkpass/internal/util"
)

var cfg = config.LoadConfig()

func HandleSendPhoto(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var headers = dto.Headers{
		Authorization: event.Headers["Authorization"],
	}
	var bucketName = cfg.ProfileBucketName

	userID, err := jwt.VerifyToken(headers.Authorization)

	if err != nil {
		return util.EncodeJson(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	storage := storage.NewStorage()
	s := service.NewSendPhotoService(storage)

	requestID, _ := ctx.Value(contextkeys.RequestID).(string)

	logger.Log.Info("send user photo",
		"request_id", requestID,
	)

	url, err := s.SendPhoto(ctx, userID, bucketName)

	if err != nil {
		return util.EncodeJson(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return util.EncodeJson(http.StatusOK, map[string]any{
		"photoUrl":  url.PhotoUrl,
		"uploadUrl": url.UploadUrl,
	})
}
