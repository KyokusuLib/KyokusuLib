package service

import (
	"context"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"github.com/lanxre/kyokusulib/internal/sse"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
	Hub  *sse.NotificationHub
}

func NewNotificationService(repo *repository.NotificationRepository, hub *sse.NotificationHub) *NotificationService {
	return &NotificationService{Repo: repo, Hub: hub}
}

func (s *NotificationService) Create(ctx context.Context, userID int64, title, message string) (*dto.Notification, error) {
	model := &db.Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
	}
	saved, err := s.Repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	notification := &dto.Notification{
		ID:        saved.ID,
		UserID:    saved.UserID,
		Title:     saved.Title,
		Message:   saved.Message,
		IsRead:    saved.IsRead,
		CreatedAt: saved.CreatedAt,
	}

	s.Hub.Publish(ctx, userID, *notification)

	return notification, nil
}

func (s *NotificationService) CreateBatch(ctx context.Context, userIDs []int64, title, message string) error {
	if len(userIDs) == 0 {
		return nil
	}

	saved, err := s.Repo.CreateBatch(ctx, userIDs, title, message)
	if err != nil {
		return err
	}

	for _, n := range saved {
		s.Hub.Publish(ctx, n.UserID, dto.Notification{
			ID:        n.ID,
			UserID:    n.UserID,
			Title:     n.Title,
			Message:   n.Message,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		})
	}
	return nil
}

func (s *NotificationService) GetByUserID(ctx context.Context, userID int64, params dto.QueryParams) ([]*dto.Notification, error) {
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 50
	}

	notifications, err := s.Repo.GetByUserID(ctx, userID, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.Notification, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, &dto.Notification{
			ID:        n.ID,
			UserID:    n.UserID,
			Title:     n.Title,
			Message:   n.Message,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		})
	}
	return result, nil
}

func (s *NotificationService) GetSinceID(ctx context.Context, userID int64, lastID int64) ([]*dto.Notification, error) {
	notifications, err := s.Repo.GetSinceID(ctx, userID, lastID, 50)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.Notification, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, &dto.Notification{
			ID:        n.ID,
			UserID:    n.UserID,
			Title:     n.Title,
			Message:   n.Message,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		})
	}
	return result, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	return s.Repo.GetUnreadCount(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	return s.Repo.MarkAsRead(ctx, notificationID, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.Repo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) Delete(ctx context.Context, id, userID int64) error {
	return s.Repo.Delete(ctx, id, userID)
}

func (s *NotificationService) GetNotificationStats(ctx context.Context, userID int64) (*dto.NotificationStats, error) {
	unreadCount, err := s.Repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return s.mapNotificatonStats(int64(unreadCount)), nil
}

func (s *NotificationService) mapNotificatonStats(unreadCount int64) *dto.NotificationStats {
	return &dto.NotificationStats{
		UnreadCount: unreadCount,
	}
}