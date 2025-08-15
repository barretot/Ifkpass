package dto

type VerifyEmailInput struct {
	Email    string `json:"email" validate:"required,email"`
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}
