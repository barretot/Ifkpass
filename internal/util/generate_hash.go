package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateSecretHash(clientSecret, attr, clientID string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(attr + clientID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
