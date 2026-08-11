package routes

import (
	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/middleware"
)

func (a *UserExperianceRoutes) Register(cfg *config.Config, r *mux.Router) {
	s := r.PathPrefix("/api/user/experiance").Subrouter()
	
	s.HandleFunc("/definitions", middleware.AuthMiddleware(middleware.RoleGuard(a.Handler.GetDefinitions, middleware.RoleModerator), cfg.JWTSecret)).Methods("GET")
	s.HandleFunc("", middleware.AuthMiddleware(middleware.RoleGuard(a.Handler.UpdateUserLevel, middleware.RoleModerator), cfg.JWTSecret)).Methods("PUT")
}