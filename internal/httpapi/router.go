package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/digocelo/account-api/internal/account"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(log *slog.Logger, svc *account.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(RequestID)
	r.Use(Logger(log))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	h := NewHandlers(svc)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/accounts", h.CreateAccount)
		r.Get("/accounts/{id}", h.GetAccount)
		r.Put("/accounts/{id}", h.UpdateAccount)
	})
	return r
}
