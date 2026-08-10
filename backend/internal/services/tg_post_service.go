package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"github.com/lanxre/kyokusulib/internal/sse"
	"github.com/lanxre/kyokusulib/internal/utils/files"
)

const tgDeleteTimeout = 10 * time.Second

type TgPostService struct {
	Repo *repository.TgPostRepository
	Hub  *sse.TgPostsHub
	cfg  *config.Config
	log  *slog.Logger
}

func NewTgPostService(repo *repository.TgPostRepository, hub *sse.TgPostsHub, cfg *config.Config, log *slog.Logger) *TgPostService {
	return &TgPostService{Repo: repo, Hub: hub, cfg: cfg, log: log}
}

func (s *TgPostService) Create(ctx context.Context, post *db.TgPost) (*dto.TgPost, error) {
	saved, err := s.Repo.Create(ctx, post)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, nil
	}

	result := toTgPostDTO(saved)

	s.Hub.Publish(ctx, *result)

	return result, nil
}

func (s *TgPostService) List(ctx context.Context, limit, offset int) ([]*dto.TgPost, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	posts, err := s.Repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.TgPost, 0, len(posts))
	for _, p := range posts {
		result = append(result, toTgPostDTO(p))
	}
	return result, nil
}

func (s *TgPostService) Delete(ctx context.Context, id int64) error {
	post, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	for _, img := range post.Images {
		if err := files.DeleteImage(img.ImagePath); err != nil {
			return err
		}
	}

	if err := s.Repo.Delete(ctx, id); err != nil {
		return err
	}

	s.Hub.PublishDelete(ctx, id)
	s.deleteFromTelegram(post.MessageID)
	return nil
}

func (s *TgPostService) GetSinceID(ctx context.Context, lastID int64, limit int) ([]*dto.TgPost, error) {
	posts, err := s.Repo.GetSinceID(ctx, lastID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.TgPost, 0, len(posts))
	for _, p := range posts {
		result = append(result, toTgPostDTO(p))
	}
	return result, nil
}

func (s *TgPostService) deleteFromTelegram(messageID int64) {
	if s.cfg.TgBotToken == "" || s.cfg.TgChannelUsername == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), tgDeleteTimeout)
	defer cancel()

	endpoint := fmt.Sprintf(
		"https://api.telegram.org/bot%s/deleteMessage",
		s.cfg.TgBotToken,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(
		url.Values{
			"chat_id":    {"@" + s.cfg.TgChannelUsername},
			"message_id": {strconv.FormatInt(messageID, 10)},
		}.Encode(),
	))
	if err != nil {
		s.log.Warn(
			"failed to build telegram delete request",
			slog.Int64("message_id", messageID),
			slog.String("error", err.Error()),
		)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.log.Warn(
			"failed to delete telegram post",
			slog.Int64("message_id", messageID),
			slog.String("error", err.Error()),
		)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		s.log.Warn(
			"telegram rejected delete request",
			slog.Int64("message_id", messageID),
			slog.Int("status", resp.StatusCode),
			slog.String("body", string(body)),
		)
	}
}

func toTgPostDTO(p *db.TgPost) *dto.TgPost {
	imageURLs := make([]string, 0, len(p.Images))
	for _, img := range p.Images {
		imageURLs = append(imageURLs, img.ImagePath)
	}

	return &dto.TgPost{
		ID:        p.ID,
		MessageID: p.MessageID,
		Text:      p.Text,
		ImageURLs: imageURLs,
		CreatedAt: p.CreatedAt,
	}
}
