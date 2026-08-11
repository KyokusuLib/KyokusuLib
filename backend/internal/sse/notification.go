package sse

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/redis/go-redis/v9"
)

type Client chan dto.Notification

type NotificationHub struct {
	mu      sync.RWMutex
	clients map[int64]map[Client]struct{}
	log     *slog.Logger
	redis   *redis.Client
}

const (
	NotificationsChannel = "notifications:stream"
)

type HubMessage struct {
	UserID       int64            `json:"user_id"`
	Notification dto.Notification `json:"notification"`
}

func NewNotificationHub(log *slog.Logger, redis *redis.Client) *NotificationHub {
	return &NotificationHub{
		clients: make(map[int64]map[Client]struct{}),
		log:     log,
		redis:   redis,
	}
}

func (h *NotificationHub) Start(ctx context.Context) {
	if h.redis == nil {
		h.log.Warn("redis not available — notification hub running in local-only mode")
		return
	}

	pubsub := h.redis.Subscribe(ctx, NotificationsChannel)

	go func() {
		defer pubsub.Close()

		ch := pubsub.Channel()
		for msg := range ch {
			var hubMsg HubMessage
			if err := json.Unmarshal([]byte(msg.Payload), &hubMsg); err != nil {
				h.log.Error("failed to unmarshal notification hub message", slog.String("error", err.Error()))
				continue
			}

			h.publishToLocalClients(hubMsg.UserID, hubMsg.Notification)
		}
	}()
}

func (h *NotificationHub) Subscribe(userID int64) Client {
	h.mu.Lock()
	defer h.mu.Unlock()

	client := make(Client, 64)

	if _, exists := h.clients[userID]; !exists {
		h.clients[userID] = make(map[Client]struct{})
	}

	h.clients[userID][client] = struct{}{}

	return client
}

func (h *NotificationHub) Unsubscribe(userID int64, client Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userClients, ok := h.clients[userID]
	if !ok {
		return
	}

	delete(userClients, client)
	close(client)

	if len(userClients) == 0 {
		delete(h.clients, userID)
	}
}

func (h *NotificationHub) Publish(
	ctx context.Context,
	userID int64,
	notification dto.Notification,
) {
	if h.redis == nil {
		h.publishToLocalClients(userID, notification)
		return
	}

	hubMsg := HubMessage{
		UserID:       userID,
		Notification: notification,
	}

	data, err := json.Marshal(hubMsg)
	if err != nil {
		h.log.Error("failed to marshal notification hub message", slog.String("error", err.Error()))
		return
	}

	if err := h.redis.Publish(ctx, NotificationsChannel, data).Err(); err != nil {
		h.log.Error("failed to publish notification to redis", slog.String("error", err.Error()))
		h.publishToLocalClients(userID, notification)
	}
}

func (h *NotificationHub) publishToLocalClients(
	userID int64,
	notification dto.Notification,
) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[userID]
	if !ok {
		return
	}

	for client := range clients {
		select {
		case client <- notification:
		default:
			h.log.Warn("notification buffer overflow",
				slog.Int64("user_id", userID),
				slog.Int64("notification_id", notification.ID),
			)
		}
	}
}

func (h *NotificationHub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, clients := range h.clients {
		for client := range clients {
			close(client)
		}
	}
	h.clients = make(map[int64]map[Client]struct{})
}

func (h *NotificationHub) PublishMany(
	ctx context.Context,
	userIDs []int64,
	notification dto.Notification,
) {
	for _, userID := range userIDs {
		h.Publish(ctx, userID, notification)
	}
}
