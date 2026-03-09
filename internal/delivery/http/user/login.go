package user

import (
	"errors"
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	out, err := h.srv.Login(r.Context(), user.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})
}
