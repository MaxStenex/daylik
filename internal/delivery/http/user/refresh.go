package user

import (
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	api "github.com/maximrozinkevich/daylik/internal/generated/api"
	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	token := refreshTokenFromRequest(r)
	if token == "" {
		var req api.RefreshRequest
		if !httputil.BindJSON(w, r, &req) {
			return
		}
		token = *req.RefreshToken
	}

	if token == "" {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("missing refresh token"))
		return
	}

	out, err := h.srv.Refresh(r.Context(), user.RefreshInput{
		RefreshToken: token,
	})
	if err != nil {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp(err.Error()))
		return
	}

	h.setRefreshCookie(w, out.RefreshToken)

	httputil.WriteJSON(w, http.StatusOK, api.TokenResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})
}
