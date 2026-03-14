package user

import (
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	api "github.com/maximrozinkevich/daylik/internal/generated/api"
	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("unauthorized"))
		return
	}

	var req api.LogoutRequest
	if !httputil.BindJSON(w, r, &req) {
		return
	}

	if err := h.srv.Logout(r.Context(), user.LogoutInput{
		UserID:       userID,
		RefreshToken: req.RefreshToken,
	}); err != nil {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
