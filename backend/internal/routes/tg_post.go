package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/handlers"
	"github.com/lanxre/kyokusulib/internal/middleware"
)

type TgPostRoutes struct {
	Handler *handlers.TgPostHandler
}

func (t *TgPostRoutes) Register(cfg *config.Config, r *mux.Router) {
	s := r.PathPrefix("/api/tg").Subrouter()

	s.Handle("/posts", http.HandlerFunc(t.Handler.Ingest)).Methods("POST")
	s.Handle("/posts", http.HandlerFunc(t.Handler.List)).Methods("GET")
	s.HandleFunc("/posts/{id:[0-9]+}",
		middleware.AuthMiddleware(middleware.RoleGuard(t.Handler.Delete, middleware.RoleModerator), cfg.JWTSecret),
	).Methods("DELETE")
	s.Handle("/stream", http.HandlerFunc(t.Handler.Stream)).Methods("GET")
}
