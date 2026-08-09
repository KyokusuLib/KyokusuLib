package routes

import (
	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/middleware"
)

func (r *ModerationRoutes) Register(cfg *config.Config, router *mux.Router) {
	sub := router.PathPrefix("/api/moderation").Subrouter()

	sub.Handle("/pending", middleware.AuthMiddleware(
		middleware.RoleGuard(r.Handler.GetPending, middleware.RoleModerator),
		cfg.JWTSecret,
	)).Methods("GET")

	sub.Handle("/volumes/{id:[a-fA-F0-9-]+}/approve", middleware.AuthMiddleware(
		middleware.RoleGuard(r.Handler.ApproveVolume, middleware.RoleModerator),
		cfg.JWTSecret,
	)).Methods("POST")

	sub.Handle("/volumes/{id:[a-fA-F0-9-]+}/reject", middleware.AuthMiddleware(
		middleware.RoleGuard(r.Handler.RejectVolume, middleware.RoleModerator),
		cfg.JWTSecret,
	)).Methods("POST")

	sub.Handle("/chapters/{id:[a-fA-F0-9-]+}/approve", middleware.AuthMiddleware(
		middleware.RoleGuard(r.Handler.ApproveChapter, middleware.RoleModerator),
		cfg.JWTSecret,
	)).Methods("POST")

	sub.Handle("/chapters/{id:[a-fA-F0-9-]+}/reject", middleware.AuthMiddleware(
		middleware.RoleGuard(r.Handler.RejectChapter, middleware.RoleModerator),
		cfg.JWTSecret,
	)).Methods("POST")

	sub.Handle("/volumes/{id:[a-fA-F0-9-]+}", middleware.AuthMiddleware(
		middleware.RoleGuard(r.Handler.UpdateVolume, middleware.RoleModerator),
		cfg.JWTSecret,
	)).Methods("PUT")

	sub.Handle("/chapters/{id:[a-fA-F0-9-]+}", middleware.AuthMiddleware(
		middleware.RoleGuard(r.Handler.UpdateChapter, middleware.RoleModerator),
		cfg.JWTSecret,
	)).Methods("PUT")
}
