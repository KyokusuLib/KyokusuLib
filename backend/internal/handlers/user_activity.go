package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/middleware"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	service "github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/utils/response"
)

type UserActivityHandler struct {
	Service *service.UserActivityService
}

func NewUserActivityHandler(userActivityService *service.UserActivityService) *UserActivityHandler {
	return &UserActivityHandler{
		Service: userActivityService,
	}
}

func (h *UserActivityHandler) CreateUserActivity(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "User not found")
		return
	}

	var userActivity dto.CreateUserActivity
	if err := json.NewDecoder(r.Body).Decode(&userActivity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateUserActivity(r.Context(), userID, &userActivity); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Activity created successfully"})
}

func (h *UserActivityHandler) GetUserActivities(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "User not found")
		return
	}

	limit, after, err := parseFeedPagination(r, 20)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	page, err := h.Service.GetUserActivity(r.Context(), userID, limit, after)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	setCursorHeaders(w, page)
	response.Success(w, http.StatusOK, page.Items)
}

func (h *UserActivityHandler) GetUserActivityById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		response.Error(w, http.StatusBadRequest, "User ID is missing")
		return
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	limit, after, err := parseFeedPagination(r, 20)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	page, err := h.Service.GetUserActivity(r.Context(), userID, limit, after)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	setCursorHeaders(w, page)
	response.Success(w, http.StatusOK, page.Items)
}
