package habit

import (
	"net/http"

	api "github.com/maximrozinkevich/daylik/internal/generated/api"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("Unauthorized"))
		return
	}

	out, err := h.srv.List(r.Context(), usecase.ListInput{UserID: userID})
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrResp("Internal server error"))
		return
	}

	habits := make([]api.HabitResponse, 0, len(out.Habits))
	for _, item := range out.Habits {
		var todayLog *api.HabitLogResponse
		if item.TodayLog != nil {
			todayLog = &api.HabitLogResponse{
				Id:             item.TodayLog.ID,
				HabitId:        item.TodayLog.HabitID,
				CompletedCount: item.TodayLog.CompletedCount,
				CreatedAt:      item.TodayLog.CreatedAt,
			}
		}

		habits = append(habits, api.HabitResponse{
			Id:          item.ID,
			Name:        item.Name,
			ExpReward:   item.ExpReward,
			DailyTarget: item.DailyTarget,
			Unit:        item.Unit,
			CreatedAt:   item.CreatedAt,
			TodayLog:    todayLog,
		})
	}

	httputil.WriteJSON(w, http.StatusOK, api.ListHabitsResponse{Habits: habits})
}
