package habits_log

import (
	"net/http"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	"github.com/maximrozinkevich/daylik/internal/generated/api"
	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit_log"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("Unauthorized"))
		return
	}

	var req api.CreateHabitLogItemJSONRequestBody
	if !httputil.BindJSON(w, r, &req) {
		return
	}

	out, err := h.srv.Create(ctx, usecase.CreateInput{
		UserID:         userID,
		HabitID:        req.HabitId,
		CompletedCount: req.CompletedCount,
	})
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrResp(err.Error()))
		return
	}

	habitLog := out.HabitLog

	httputil.WriteJSON(w, http.StatusCreated, api.HabitLogResponse{
		Id:             habitLog.ID,
		HabitId:        habitLog.HabitID,
		CompletedCount: habitLog.CompletedCount,
		CreatedAt:      habitLog.CreatedAt,
	})
}
