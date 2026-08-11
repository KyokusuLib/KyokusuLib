package tg

import (
	"context"
	"errors"

	"github.com/lanxre/kyokusulib/internal/apperrors"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"github.com/lanxre/kyokusulib/internal/sse"
	"github.com/lanxre/kyokusulib/internal/utils/files"
)

type TgPostService struct {
	Repo *repository.TgPostRepository
	Hub  *sse.TgPostsHub
	Bot  *TgBotAPI
}

func NewTgPostService(repo *repository.TgPostRepository, hub *sse.TgPostsHub, botApi *TgBotAPI) *TgPostService {
	return &TgPostService{Repo: repo, Hub: hub, Bot: botApi}
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

	if err := s.Bot.DeleteFromTelegram(ctx, post.MessageID); err != nil && !errors.Is(err, apperrors.ErrTelegramDeletePermanent) {
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
