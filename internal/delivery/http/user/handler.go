package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

type service interface {
	Register(ctx context.Context, in user.RegisterInput) (user.RegisterOutput, error)
	Login(ctx context.Context, in user.LoginInput) (user.LoginOutput, error)
	Refresh(ctx context.Context, in user.RefreshInput) (user.RefreshOutput, error)
	Logout(ctx context.Context, in user.LogoutInput) error
}

type Handler struct {
	srv service
}

func NewHandler(srv service) *Handler {
	return &Handler{srv: srv}
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
