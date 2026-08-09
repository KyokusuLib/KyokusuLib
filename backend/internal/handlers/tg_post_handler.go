package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/models/db"
	service "github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/sse"
	"github.com/lanxre/kyokusulib/internal/utils/conv"
	"github.com/lanxre/kyokusulib/internal/utils/files"
	"github.com/lanxre/kyokusulib/internal/utils/response"
)

type TgPostHandler struct {
	Service *service.TgPostService
	Hub     *sse.TgPostsHub
	cfg     *config.Config
}

func NewTgPostHandler(tgPostService *service.TgPostService, hub *sse.TgPostsHub, cfg *config.Config) *TgPostHandler {
	return &TgPostHandler{
		Service: tgPostService,
		Hub:     hub,
		cfg:     cfg,
	}
}

func (h *TgPostHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	if h.cfg.TgBotApiKey == "" || r.Header.Get("X-API-Key") != h.cfg.TgBotApiKey {
		response.Error(w, http.StatusForbidden, "Invalid API key")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	messageID, err := strconv.ParseInt(r.FormValue("message_id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid message_id")
		return
	}

	text := r.FormValue("text")

	headers := r.MultipartForm.File["images"]
	images := make([]db.TgPostImage, 0, len(headers))
	for i, header := range headers {
		file, err := header.Open()
		if err != nil {
			h.cleanupImages(images)
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		imageURL, err := files.UploadImage(r.Context(), file, header, "tg", 0, 0)
		file.Close()
		if err != nil {
			h.cleanupImages(images)
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		images = append(images, db.TgPostImage{Position: i, ImagePath: imageURL})
	}

	post, err := h.Service.Create(r.Context(), &db.TgPost{
		MessageID: messageID,
		Text:      text,
		Images:    images,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if post == nil {
		response.SuccessOkEmpty(w)
		return
	}

	response.SuccessWithEntity(w, http.StatusCreated, post)
}

func (h *TgPostHandler) cleanupImages(images []db.TgPostImage) {
	for _, img := range images {
		_ = files.DeleteImage(img.ImagePath)
	}
}

func (h *TgPostHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit := conv.StringToInt(query.Get("limit"))
	offset := conv.StringToInt(query.Get("offset"))

	if limit < 0 {
		response.Error(w, http.StatusBadRequest, "limit must be non-negative")
		return
	}
	if offset < 0 {
		response.Error(w, http.StatusBadRequest, "offset must be non-negative")
		return
	}
	if limit > 100 {
		response.Error(w, http.StatusBadRequest, "limit must not exceed 100")
		return
	}

	posts, err := h.Service.List(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithEntity(w, http.StatusOK, posts)
}

func (h *TgPostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	err = h.Service.Delete(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		response.Error(w, http.StatusNotFound, "Post not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *TgPostHandler) Stream(w http.ResponseWriter, r *http.Request) {
	sw, ok := sse.NewSSEWriter(w)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	if lastIDStr := r.Header.Get("Last-Event-ID"); lastIDStr != "" {
		lastID, err := strconv.ParseInt(lastIDStr, 10, 64)
		if err == nil && lastID > 0 {
			posts, err := h.Service.GetSinceID(r.Context(), lastID, 50)
			if err == nil {
				for _, p := range posts {
					data, err := json.Marshal(p)
					if err != nil {
						continue
					}
					if err := sw.SendEvent("tg_post", p.ID, data); err != nil {
						return
					}
				}
			}
		}
	}

	client := h.Hub.Subscribe()
	defer h.Hub.Unsubscribe(client)

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
		case ev, ok := <-client:
			if !ok {
				return
			}
			switch ev.Type {
			case sse.TgPostEventPost:
				if ev.Post == nil {
					continue
				}
				data, err := json.Marshal(ev.Post)
				if err != nil {
					continue
				}
				if err := sw.SendEvent("tg_post", ev.Post.ID, data); err != nil {
					return
				}
			case sse.TgPostEventDelete:
				data, err := json.Marshal(struct {
					ID int64 `json:"id"`
				}{ev.ID})
				if err != nil {
					continue
				}
				if err := sw.SendEvent("tg_post_deleted", ev.ID, data); err != nil {
					return
				}
			}
		}
	}
}
