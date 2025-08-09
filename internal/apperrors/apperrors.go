package apperrors

import "errors"

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInternalServerError = errors.New("internal server error")
	ErrUserNotFound        = errors.New("user not found")
	ErrFailedToGetUserId   = errors.New("failed to get user id")
)
