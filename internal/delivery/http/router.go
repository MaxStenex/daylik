package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/maximrozinkevich/daylik/internal/delivery/http/habit"
	habits_log "github.com/maximrozinkevich/daylik/internal/delivery/http/habit_log"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/middleware"
	"github.com/maximrozinkevich/daylik/internal/delivery/http/user"
)

func NewRouter(
	log *slog.Logger,
	tokens middleware.TokenVerifier,
	userHandler *user.Handler,
	habitHandler *habit.Handler,
	habitLogHandler *habits_log.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.SlogLogger(log))
	r.Use(chimiddleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
			r.Post("/refresh", userHandler.Refresh)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(tokens))

			r.Post("/auth/logout", userHandler.Logout)

			r.Route("/habits", func(r chi.Router) {
				r.Post("/", habitHandler.Create)
				r.Get("/", habitHandler.List)
				r.Put("/{id}", habitHandler.Update)
				r.Delete("/{id}", habitHandler.Delete)
			})

			r.Route("/habits-log", func(r chi.Router) {
				r.Post("/", habitLogHandler.Create)
				r.Get("/", habitLogHandler.List)
				r.Put("/{id}", habitLogHandler.Update)
			})
		})
	})

	return r
}
