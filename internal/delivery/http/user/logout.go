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
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("Unauthorized"))
		return
	}

	refreshToken := refreshTokenFromRequest(r)
	if refreshToken == "" {
		var req api.RefreshRequest
		if !httputil.BindJSON(w, r, &req) {
			return
		}
		refreshToken = *req.RefreshToken
	}

	if refreshToken == "" {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("Missing refresh token"))
		return
	}

	if err := h.srv.Logout(r.Context(), user.LogoutInput{
		UserID:       userID,
		RefreshToken: refreshToken,
	}); err != nil {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp(err.Error()))
		return
	}

	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
