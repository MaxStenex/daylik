package habit

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	api "github.com/maximrozinkevich/daylik/internal/generated/api"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("Unauthorized"))
		return
	}

	habitID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrResp("Invalid habit ID"))
		return
	}

	var req api.UpdateHabitJSONRequestBody
	if !httputil.BindJSON(w, r, &req) {
		return
	}

	err = h.srv.Update(r.Context(), usecase.UpdateInput{
		HabitID:     habitID,
		UserID:      userID,
		Name:        req.Name,
		ExpReward:   req.ExpReward,
		DailyTarget: req.DailyTarget,
		Unit:        req.Unit,
	})
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrResp(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
