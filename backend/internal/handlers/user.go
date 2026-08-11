package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/middleware"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	service "github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/utils/response"
)

type UserHandler struct {
	UserService 	 	*service.UserService
	Validator    	 	*validator.Validate
	ActivityService  	*service.UserActivityService
	NotificationService *service.NotificationService
}

func NewUserHandler(userService *service.UserService, validator *validator.Validate, activityService *service.UserActivityService, notificationService *service.NotificationService) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Validator: validator,
		ActivityService: activityService,
		NotificationService: notificationService,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	limit := 20
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		offset = o
	}

	users, err := h.UserService.GetUsers(r.Context(), search, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get users")
		return
	}

	response.SuccessWithEntity(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.UserService.GetUserById(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if user == nil {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	moderatorId, _ := r.Context().Value(middleware.UserIDKey).(int)
	userID, _ := strconv.Atoi(mux.Vars(r)["userId"])
	
	var req dto.UpdateUserStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Print(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	if err := h.UserService.UpdateUserStatus(r.Context(), userID, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
	activityMetaData, _ := json.Marshal(map[string]any{
   	"user_id": moderatorId,
   	"title": "Произведенно обновление статуса аккаунт пользователя",
   	"status": req.Status,
    "time": time.Now().Format("02-01-2006"),
	})
	
	h.ActivityService.CreateUserActivity(r.Context(), moderatorId, &dto.CreateUserActivity{
		ActivityType: "update_account_status",
		TargetID: userID,
		Metadata: activityMetaData,
	})
	
	h.NotificationService.Create(r.Context(), int64(moderatorId), "Обновление статуса аккаунта", fmt.Sprintf("Обновил статус аккаунта пользователя: %d", userID))


	response.Success(w, http.StatusOK, "User status was updated: "+time.Now().Format("02-01-2006"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	moderatorId, _ := r.Context().Value(middleware.UserIDKey).(int)
	userID, _ := strconv.Atoi(mux.Vars(r)["userId"])

	var req dto.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	if err := h.UserService.UpdateUser(r.Context(), userID, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	activityMetaData, _ := json.Marshal(map[string]any{
		"user_id": moderatorId,
		"title":   "Произведено обновление данных аккаунта пользователя",
		"time":    time.Now().Format("02-01-2006"),
	})

	h.ActivityService.CreateUserActivity(r.Context(), moderatorId, &dto.CreateUserActivity{
		ActivityType: "update_account",
		TargetID:     userID,
		Metadata:     activityMetaData,
	})

	h.NotificationService.Create(r.Context(), int64(moderatorId), "Обновление данных аккаунта", fmt.Sprintf("Обновил данные аккаунта пользователя: %d", userID))

	response.Success(w, http.StatusOK, "User was updated: "+time.Now().Format("02-01-2006"))
}

func (h *UserHandler) UpdateMyStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	var req dto.UpdateUserStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.UserService.UpdateUserStatus(r.Context(), userID, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(w, http.StatusOK, "OK")
}

func (h *UserHandler) UpdateUserTag(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}
	
	var req dto.UpdateUserTagDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.UserService.UpdateUserTag(r.Context(), userID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Тэг профиля успешно обновлен")
}

func (h *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	moderatorId, _ := r.Context().Value(middleware.UserIDKey).(int)
	userID, _ := strconv.Atoi(mux.Vars(r)["userId"])

	if err := h.UserService.DeleteUserById(r.Context(), userID); err != nil {
		response.Error(w, http.StatusBadRequest, "Unexpected error: "+err.Error())
		return
	}

	activityMetaData, _ := json.Marshal(map[string]any{
    	"user_id": moderatorId,
     	"title": "Произведенно удаление аккаунта пользователя",
     	"delete_account_id": userID,
        "time": time.Now().Format("02-01-2006"),
	})
	
	h.ActivityService.CreateUserActivity(r.Context(), moderatorId, &dto.CreateUserActivity{
		ActivityType: "delete_user",
		TargetID: userID,
		Metadata: activityMetaData,
	})

	h.NotificationService.Create(r.Context(), int64(moderatorId), "Удаление аккаунта", fmt.Sprintf("Удалил аккаунта пользователя: %d", userID))

	response.Success(w, http.StatusOK, "User was delete on: "+time.Now().Format("02-01-2006"))
}