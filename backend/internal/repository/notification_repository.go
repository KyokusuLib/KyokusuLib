package repository

import (
	"context"
	"database/sql"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lib/pq"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *db.Notification) (*db.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, title, message, is_read, created_at)
		VALUES ($1, $2, $3, FALSE, NOW())
		RETURNING id, created_at`
	err := r.DB.QueryRowContext(ctx, query, n.UserID, n.Title, n.Message).
		Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *NotificationRepository) CreateBatch(ctx context.Context, userIDs []int64, title, message string) ([]*db.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, title, message, is_read, created_at)
		SELECT unnest($1::bigint[]), $2, $3, FALSE, NOW()
		RETURNING id, user_id, created_at`
	rows, err := r.DB.QueryContext(ctx, query, pq.Array(userIDs), title, message)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]*db.Notification, 0, len(userIDs))
	for rows.Next() {
		var n db.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.CreatedAt); err != nil {
			return nil, err
		}
		n.Title = title
		n.Message = message
		notifications = append(notifications, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*db.Notification, error) {
	query := `
		SELECT id, user_id, title, message, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*db.Notification
	for rows.Next() {
		var n db.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	if notifications == nil {
		notifications = []*db.Notification{}
	}
	return notifications, nil
}

func (r *NotificationRepository) GetSinceID(ctx context.Context, userID int64, lastID int64, limit int) ([]*db.Notification, error) {
	query := `
		SELECT id, user_id, title, message, is_read, created_at
		FROM notifications
		WHERE user_id = $1 AND id > $2
		ORDER BY id ASC
		LIMIT $3`
	rows, err := r.DB.QueryContext(ctx, query, userID, lastID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*db.Notification
	for rows.Next() {
		var n db.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	if notifications == nil {
		notifications = []*db.Notification{}
	}
	return notifications, nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`
	var count int
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1 AND user_id = $2`
	_, err := r.DB.ExecContext(ctx, query, notificationID, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`
	_, err := r.DB.ExecContext(ctx, query, userID)
	return err
}

func (r *NotificationRepository) Delete(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM notifications WHERE id = $1 AND user_id = $2`
	result, err := r.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
