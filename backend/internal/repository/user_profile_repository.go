package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lib/pq"
)

type UserProfileRepository struct {
	DB *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *UserProfileRepository {
	return &UserProfileRepository{DB: db}
}

func (r *UserProfileRepository) GetUserLevel(ctx context.Context, userID int) (*db.UserLevel, error) {
	u := &db.UserLevel{}

	var (
		level          int
		experience     int64
		levelTitle     sql.NullString
		xpForThisLevel sql.NullInt64
		xpForNextLevel sql.NullInt64
	)

	query := `
	SELECT 
		p.level, p.experience,
		ld.title AS level_title,
		ld.total_xp_required AS xp_for_this_level,
		COALESCE(next_ld.total_xp_required, ld.total_xp_required) AS xp_for_next_level
	FROM user_profiles p
	LEFT JOIN level_definitions ld ON ld.level = p.level
	LEFT JOIN level_definitions next_ld ON next_ld.level = p.level + 1
	WHERE p.user_id = $1
	`

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&level, &experience, &levelTitle, &xpForThisLevel, &xpForNextLevel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	if levelTitle.Valid {
		u.LevelTitle = levelTitle.String
	}

	if xpForNextLevel.Valid && xpForThisLevel.Valid {
		needed := max(xpForNextLevel.Int64-experience, 0)
		u.XPForNext = needed
	}

	u.Level = level
	u.Experience = experience

	return u, nil
}

func (r *UserProfileRepository) GetUserLevelsBatch(ctx context.Context, userIDs []int) (map[int]*db.UserLevel, error) {
	result := make(map[int]*db.UserLevel, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT
			p.user_id, p.level, p.experience,
			ld.title AS level_title,
			ld.total_xp_required AS xp_for_this_level,
			COALESCE(next_ld.total_xp_required, ld.total_xp_required) AS xp_for_next_level
		FROM user_profiles p
		LEFT JOIN level_definitions ld ON ld.level = p.level
		LEFT JOIN level_definitions next_ld ON next_ld.level = p.level + 1
		WHERE p.user_id = ANY($1)
	`

	rows, err := r.DB.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		u := &db.UserLevel{}
		var (
			levelTitle     sql.NullString
			xpForThisLevel sql.NullInt64
			xpForNextLevel sql.NullInt64
		)

		if err := rows.Scan(&userID, &u.Level, &u.Experience, &levelTitle, &xpForThisLevel, &xpForNextLevel); err != nil {
			return nil, err
		}

		if levelTitle.Valid {
			u.LevelTitle = levelTitle.String
		}

		if xpForNextLevel.Valid && xpForThisLevel.Valid {
			u.XPForNext = max(xpForNextLevel.Int64-u.Experience, 0)
		}

		result[userID] = u
	}

	return result, rows.Err()
}

func (r *UserProfileRepository) GetLevelDefinition(ctx context.Context, level int) (*db.LevelDefinition, error) {
	query := `
		SELECT * FROM level_definitions
		WHERE level = $1
	`

	var ld db.LevelDefinition
	err := r.DB.QueryRowContext(ctx, query, level).Scan(&ld.Level, &ld.Title, &ld.TotalXpRequired)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &ld, nil
}

func (r *UserProfileRepository) UpUserLevel(ctx context.Context, userID int) error {
	query := `
		UPDATE user_profiles
		SET level = level + 1
		WHERE user_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, userID)
	return err
}

func (r *UserProfileRepository) SetUserLevel(ctx context.Context, userID int, level int) error {
	query := `
		UPDATE user_profiles
		SET level = $2
		WHERE user_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, userID, level)
	return err
}

func (r *UserProfileRepository) SetUserExperiance(ctx context.Context, userID int, experiance int64) error {
	query := `
		UPDATE user_profiles
		SET experience = $2
		WHERE user_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, userID, experiance)
	return err
}
