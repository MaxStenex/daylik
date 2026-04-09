package habit

import (
	"net/http"

	api "github.com/maximrozinkevich/daylik/internal/generated/api"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, httputil.ErrResp("unauthorized"))
		return
	}

	var req api.CreateHabitJSONRequestBody
	if !httputil.BindJSON(w, r, &req) {
		return
	}

	out, err := h.srv.Create(r.Context(), usecase.CreateInput{
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

	httputil.WriteJSON(w, http.StatusCreated, api.CreateHabitResponse{
		Id:          out.Habit.ID,
		Name:        out.Habit.Name,
		ExpReward:   out.Habit.ExpReward,
		DailyTarget: out.Habit.DailyTarget,
		Unit:        out.Habit.Unit,
		CreatedAt:   out.Habit.CreatedAt,
	})
}
