package routes

import (
	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/middleware"
)

func (a *UserRoutes) Register(cfg *config.Config, r *mux.Router) {
	s := r.PathPrefix("/api/user").Subrouter()

	s.HandleFunc("", middleware.AuthMiddleware(middleware.RoleGuard(a.Handler.GetUsers, middleware.RoleModerator), cfg.JWTSecret)).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}", a.Handler.GetUserById).Methods("GET")
	
	s.HandleFunc("/activity", middleware.AuthMiddleware(a.Handler.UpdateMyStatus, cfg.JWTSecret)).Methods("POST")

	s.HandleFunc("/tag", middleware.AuthMiddleware(a.Handler.UpdateUserTag, cfg.JWTSecret)).Methods("PUT")
	s.HandleFunc("/{userId:[0-9]+}/status", middleware.AuthMiddleware(middleware.RoleGuard(a.Handler.UpdateUserStatus, middleware.RoleModerator), cfg.JWTSecret)).Methods("PUT") 

	s.HandleFunc("/{userId:[0-9]+}", middleware.AuthMiddleware(middleware.RoleGuard(a.Handler.UpdateUser, middleware.RoleModerator), cfg.JWTSecret)).Methods("PUT")

	s.HandleFunc("/{userId:[0-9]+}", middleware.AuthMiddleware(middleware.RoleGuard(a.Handler.DeleteUserById, middleware.RoleAdmin), cfg.JWTSecret)).Methods("DELETE") 
}
