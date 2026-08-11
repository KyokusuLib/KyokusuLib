package sse

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/redis/go-redis/v9"
)

type TgPostEvent struct {
	Type string      `json:"type"`
	Post *dto.TgPost `json:"post,omitempty"`
	ID   int64       `json:"id,omitempty"`
}

const (
	TgPostsChannel    = "tgposts:stream"
	TgPostEventPost   = "post"
	TgPostEventDelete = "delete"
)

type TgPostsClient chan TgPostEvent

type TgPostsHub struct {
	mu      sync.RWMutex
	clients map[TgPostsClient]struct{}
	log     *slog.Logger
	redis   *redis.Client
}

func NewTgPostsHub(log *slog.Logger, redis *redis.Client) *TgPostsHub {
	return &TgPostsHub{
		clients: make(map[TgPostsClient]struct{}),
		log:     log,
		redis:   redis,
	}
}

func (h *TgPostsHub) Start(ctx context.Context) {
	if h.redis == nil {
		h.log.Warn("redis not available — tg posts hub running in local-only mode")
		return
	}

	pubsub := h.redis.Subscribe(ctx, TgPostsChannel)

	go func() {
		defer pubsub.Close()

		ch := pubsub.Channel()
		for msg := range ch {
			var ev TgPostEvent
			if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
				h.log.Error("failed to unmarshal tg post hub message", slog.String("error", err.Error()))
				continue
			}
			h.publishToLocalClients(ev)
		}
	}()
}

func (h *TgPostsHub) Subscribe() TgPostsClient {
	h.mu.Lock()
	defer h.mu.Unlock()

	client := make(TgPostsClient, 64)
	h.clients[client] = struct{}{}

	return client
}

func (h *TgPostsHub) Unsubscribe(client TgPostsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, client)
	close(client)
}

func (h *TgPostsHub) Publish(ctx context.Context, post dto.TgPost) {
	h.broadcast(ctx, TgPostEvent{Type: TgPostEventPost, Post: &post})
}

func (h *TgPostsHub) PublishDelete(ctx context.Context, id int64) {
	h.broadcast(ctx, TgPostEvent{Type: TgPostEventDelete, ID: id})
}

func (h *TgPostsHub) broadcast(ctx context.Context, ev TgPostEvent) {
	if h.redis == nil {
		h.publishToLocalClients(ev)
		return
	}

	data, err := json.Marshal(ev)
	if err != nil {
		h.log.Error("failed to marshal tg post hub message", slog.String("error", err.Error()))
		return
	}

	if err := h.redis.Publish(ctx, TgPostsChannel, data).Err(); err != nil {
		h.log.Error("failed to publish tg post to redis", slog.String("error", err.Error()))
		h.publishToLocalClients(ev)
	}
}

func (h *TgPostsHub) publishToLocalClients(ev TgPostEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client <- ev:
		default:
			h.log.Warn("tg post buffer overflow",
				slog.String("event", ev.Type),
			)
		}
	}
}

func (h *TgPostsHub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client)
	}
	h.clients = make(map[TgPostsClient]struct{})
}
