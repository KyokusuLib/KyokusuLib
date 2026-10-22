package service

import (
	"context"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"github.com/lanxre/kyokusulib/internal/utils/cursor"
)

type UserActivityService struct {
	Repo *repository.UserActivityRepository
}

func NewUserActivityService(repo *repository.UserActivityRepository) *UserActivityService {
	return &UserActivityService{Repo: repo}
}

func (s *UserActivityService) GetUserActivity(ctx context.Context, userID int, limit int, after *cursor.Cursor) (dto.CursorPage[*dto.GetUserActivity], error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	activities, hasMore, err := s.Repo.GetUserActivies(ctx, userID, limit, after)
	if err != nil {
		return dto.CursorPage[*dto.GetUserActivity]{}, err
	}

	userActivity := make([]*dto.GetUserActivity, 0, len(activities))
	for _, activity := range activities {
		userActivity = append(userActivity, &dto.GetUserActivity{
			ID:           activity.ID,
			ActivityType: activity.ActivityType,
			TargetID:     activity.TargetID,
			Timestamp:    activity.Timestamp,
			Metadata:     activity.Metadata,
		})
	}

	page := dto.CursorPage[*dto.GetUserActivity]{Items: userActivity, HasMore: hasMore}
	if hasMore && len(userActivity) > 0 {
		last := userActivity[len(userActivity)-1]
		page.NextCursor = cursor.Encode(last.Timestamp, int64(last.ID))
	}
	return page, nil
}

func (s *UserActivityService) CreateUserActivity(ctx context.Context, userID int, userActivity *dto.CreateUserActivity) error {
	model := &db.UserActivity{
		UserID:       userID,
		ActivityType: userActivity.ActivityType,
		TargetID:     userActivity.TargetID,
		Metadata:     userActivity.Metadata,
	}
	return s.Repo.CreateUserActivity(ctx, model)
}
