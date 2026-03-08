package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

type service interface {
	Register(ctx context.Context, in user.RegisterInput) (user.RegisterOutput, error)
}

type Handler struct {
	srv service
}

func NewHandler(srv service) *Handler {
	return &Handler{srv: srv}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	out, err := h.srv.Register(r.Context(), user.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail):
			writeJSON(w, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
		case errors.Is(err, user.ErrPasswordTooShort):
			writeJSON(w, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
		case errors.Is(err, user.ErrEmailTaken):
			writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, registerResponse{ID: out.ID, Email: out.Email})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
