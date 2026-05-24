package habits_log

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	api "github.com/maximrozinkevich/daylik/internal/generated/api"
	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit_log"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("unauthorized"))
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrResp("invalid habit log id"))
		return
	}

	var req api.UpdateHabitLogItemJSONRequestBody
	if !httputil.BindJSON(w, r, &req) {
		return
	}

	err = h.srv.Update(ctx, usecase.UpdateInput{
		ID:             id,
		UserID:         userID,
		CompletedCount: req.CompletedCount,
	})
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrResp(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
