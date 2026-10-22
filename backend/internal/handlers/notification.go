package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/middleware"
	service "github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/sse"
	"github.com/lanxre/kyokusulib/internal/utils/response"
)

type NotificationHandler struct {
	Service *service.NotificationService
	Hub     *sse.NotificationHub
}

func NewNotificationHandler(notificationService *service.NotificationService, hub *sse.NotificationHub) *NotificationHandler {
	return &NotificationHandler{
		Service: notificationService,
		Hub:     hub,
	}
}

func (h *NotificationHandler) StreamNotifications(w http.ResponseWriter, r *http.Request) {
	sw, ok := sse.NewSSEWriter(w)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User not found")
		return
	}

	if lastIDStr := r.Header.Get("Last-Event-ID"); lastIDStr != "" {
		lastID, err := strconv.ParseInt(lastIDStr, 10, 64)
		if err == nil && lastID > 0 {
			notifications, err := h.Service.GetSinceID(r.Context(), int64(userID), lastID)
			if err == nil {
				for _, n := range notifications {
					data, err := json.Marshal(n)
					if err != nil {
						continue
					}
					if err := sw.SendEvent("notification", n.ID, data); err != nil {
						return
					}
				}
			}
		}
	}

	client := h.Hub.Subscribe(int64(userID))
	defer h.Hub.Unsubscribe(int64(userID), client)

	if err := sw.SendHeartbeat(); err != nil {
		return
	}

	ctx := r.Context()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sw.SendHeartbeat(); err != nil {
				return
			}
		case notification, ok := <-client:
			if !ok {
				return
			}
			data, err := json.Marshal(notification)
			if err != nil {
				continue
			}
			if err := sw.SendEvent("notification", notification.ID, data); err != nil {
				return
			}
		}
	}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User not found")
		return
	}

	limit, after, err := parseFeedPagination(r, 50)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	page, err := h.Service.GetByUserID(r.Context(), int64(userID), limit, after)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	setCursorHeaders(w, page)
	response.SuccessWithEntity(w, http.StatusOK, page.Items)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User not found")
		return
	}

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		response.Error(w, http.StatusBadRequest, "Notification ID is missing")
		return
	}

	notificationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.Service.MarkAsRead(r.Context(), notificationID, int64(userID)); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User not found")
		return
	}

	if err := h.Service.MarkAllAsRead(r.Context(), int64(userID)); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User not found")
		return
	}

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		response.Error(w, http.StatusBadRequest, "Notification ID is missing")
		return
	}

	notificationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.Service.Delete(r.Context(), notificationID, int64(userID)); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NotificationHandler) NotificationStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User not found")
		return
	}

	stats, err := h.Service.GetNotificationStats(r.Context(), int64(userID))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithEntity(w, http.StatusOK, stats)
}
