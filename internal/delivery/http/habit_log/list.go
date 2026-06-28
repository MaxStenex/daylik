package habits_log

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	api "github.com/maximrozinkevich/daylik/internal/generated/api"
	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit_log"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("Unauthorized"))
		return
	}

	habitID, err := uuid.Parse(r.URL.Query().Get("habit_id"))
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrResp("Invalid habit ID"))
		return
	}

	out, err := h.srv.List(ctx, usecase.ListInput{
		UserID:  userID,
		HabitID: habitID,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidHabit), errors.Is(err, usecase.ErrForbidden):
			httputil.WriteJSON(w, http.StatusNotFound, httputil.ErrResp("Habit not found"))
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrResp("Internal server error"))
		}
		return
	}

	logs := make([]api.HabitLogResponse, 0, len(out.HabitLogs))
	for _, item := range out.HabitLogs {
		logs = append(logs, api.HabitLogResponse{
			Id:             item.ID,
			HabitId:        item.HabitID,
			CompletedCount: item.CompletedCount,
			CreatedAt:      item.CreatedAt,
		})
	}

	httputil.WriteJSON(w, http.StatusOK, api.ListHabitLogsResponse{Logs: logs})
}
