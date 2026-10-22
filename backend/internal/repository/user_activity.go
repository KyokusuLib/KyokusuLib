package repository

import (
	"context"
	"database/sql"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/utils/cursor"
)

type UserActivityRepository struct {
	DB *sql.DB
}

func NewUserActivityRepository(db *sql.DB) *UserActivityRepository {
	return &UserActivityRepository{DB: db}
}

func (r *UserActivityRepository) GetUserActivies(ctx context.Context, userID, limit int, after *cursor.Cursor) ([]*db.UserActivity, bool, error) {
	query := `
		SELECT id, activity_type, target_id, metadata, "timestamp"
		FROM user_activities
		WHERE user_id = $1
		ORDER BY "timestamp" DESC, id DESC
		LIMIT $2`
	args := []any{userID, limit + 1}
	if after != nil {
		query = `
		SELECT id, activity_type, target_id, metadata, "timestamp"
		FROM user_activities
		WHERE user_id = $1 AND ("timestamp", id) < ($2, $3)
		ORDER BY "timestamp" DESC, id DESC
		LIMIT $4`
		args = []any{userID, after.Time, after.ID, limit + 1}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	userActivity := make([]*db.UserActivity, 0, limit+1)
	for rows.Next() {
		var activity db.UserActivity
		err := rows.Scan(
			&activity.ID,
			&activity.ActivityType,
			&activity.TargetID,
			&activity.Metadata,
			&activity.Timestamp,
		)
		if err != nil {
			return nil, false, err
		}
		userActivity = append(userActivity, &activity)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(userActivity) > limit
	if hasMore {
		userActivity = userActivity[:limit]
	}
	return userActivity, hasMore, nil
}

func (r *UserActivityRepository) CreateUserActivity(ctx context.Context, userActivity *db.UserActivity) error {
	query := `
		INSERT INTO user_activities (user_id, activity_type, target_id, metadata, timestamp)
		VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.DB.ExecContext(ctx, query, userActivity.UserID, userActivity.ActivityType, userActivity.TargetID, userActivity.Metadata)
	if err != nil {
		return err
	}

	return nil
}
