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
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("unauthorized"))
		return
	}

	out, err := h.srv.List(r.Context(), usecase.ListInput{UserID: userID})
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrResp("internal server error"))
		return
	}

	habits := make([]api.HabitResponse, 0, len(out.Habits))
	for _, item := range out.Habits {
		habits = append(habits, api.HabitResponse{
			Id:          item.ID,
			Name:        item.Name,
			ExpReward:   item.ExpReward,
			DailyTarget: item.DailyTarget,
			Unit:        item.Unit,
			CreatedAt:   item.CreatedAt,
			ArchivedAt:  item.ArchivedAt,
		})
	}

	httputil.WriteJSON(w, http.StatusOK, api.ListHabitsResponse{Habits: habits})
}
