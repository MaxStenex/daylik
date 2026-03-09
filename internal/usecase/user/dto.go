package user

import "github.com/google/uuid"

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	ID    string
	Email string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string
	RefreshToken string
}

type LogoutInput struct {
	UserID       uuid.UUID
	RefreshToken string
}
