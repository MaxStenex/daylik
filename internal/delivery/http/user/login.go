package user

import (
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	api "github.com/maximrozinkevich/daylik/internal/generated/api"
	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req api.LoginRequest
	if !httputil.BindJSON(w, r, &req) {
		return
	}

	out, err := h.srv.Login(r.Context(), user.LoginInput{
		Email:    req.Email,
		Password: req.Password,
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
