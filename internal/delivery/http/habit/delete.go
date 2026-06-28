package habit

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/httputil"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = h.srv.Delete(r.Context(), usecase.DeleteInput{
		HabitID: habitID,
		UserID:  userID,
	})
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrResp(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
