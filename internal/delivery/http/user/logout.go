package user

import (
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	var req logoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if err := h.srv.Logout(r.Context(), user.LogoutInput{
		UserID:       userID,
		RefreshToken: req.RefreshToken,
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
