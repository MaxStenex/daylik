package user

import (
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	api "github.com/maximrozinkevich/daylik/internal/generated/api"
	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req api.RegisterRequest
	if !httputil.BindJSON(w, r, &req) {
		return
	}

	out, err := h.srv.Register(r.Context(), user.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp(err.Error()))
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, api.RegisterResponse{Id: out.ID, Email: out.Email})
}
