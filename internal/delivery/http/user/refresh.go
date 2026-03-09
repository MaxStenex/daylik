package user

import (
	"errors"
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	out, err := h.srv.Refresh(r.Context(), user.RefreshInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, user.ErrInvalidRefreshToken) {
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
