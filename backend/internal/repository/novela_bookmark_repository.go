package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lib/pq"
)

type NovelaBookmarkRepository struct {
	DB *sql.DB
}

func NewNovelaBookmarkRepository(db *sql.DB) *NovelaBookmarkRepository {
	return &NovelaBookmarkRepository{DB: db}
}

func (r *NovelaBookmarkRepository) SetBookmark(ctx context.Context, userID, novelaID, categoryID int) error {
	query := `
		INSERT INTO user_novela_bookmarks (user_id, novela_id, category_id, updated_at) 
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP) 
		ON CONFLICT (user_id, novela_id) 
		DO UPDATE SET 
			category_id = EXCLUDED.category_id, 
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.DB.ExecContext(ctx, query, userID, novelaID, categoryID)
	return err
}

func (r *NovelaBookmarkRepository) RemoveBookmark(ctx context.Context, userID, novelaID int) error {
	query := `DELETE FROM user_novela_bookmarks WHERE user_id = $1 AND novela_id = $2`
	_, err := r.DB.ExecContext(ctx, query, userID, novelaID)
	return err
}

func (r *NovelaBookmarkRepository) GetBookmarkStats(tx *sql.Tx, ctx context.Context, novelaID int) (*db.NovelaBookmarkSummary, error) {
	query := `
		SELECT 
			COALESCE(SUM(s.count), 0) as total_count,
			COALESCE(jsonb_object_agg(c.name, s.count), '{}'::jsonb)
		FROM (
			SELECT category_id, COUNT(*) as count 
			FROM user_novela_bookmarks 
			WHERE novela_id = $1 
			GROUP BY category_id
		) s
		JOIN bookmark_categories c ON s.category_id = c.id`

	var summary db.NovelaBookmarkSummary
	var distJSON []byte
	var totalCount int

	err := tx.QueryRowContext(ctx, query, novelaID).Scan(
		&totalCount,
		&distJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &db.NovelaBookmarkSummary{Distribution: make(map[string]int)}, nil
		}
		return nil, err
	}

	summary.TotalCount = totalCount
	summary.Distribution = make(map[string]int)
	if len(distJSON) > 0 {
		json.Unmarshal(distJSON, &summary.Distribution)
	}

	return &summary, nil
}

func (r *NovelaBookmarkRepository) GetBookmarkStatsBatch(tx *sql.Tx, ctx context.Context, novelaIDs []int) (map[int]*db.NovelaBookmarkSummary, error) {
	result := make(map[int]*db.NovelaBookmarkSummary, len(novelaIDs))
	if len(novelaIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT
			s.novela_id,
			COALESCE(SUM(s.count), 0) as total_count,
			COALESCE(jsonb_object_agg(c.name, s.count), '{}'::jsonb)
		FROM (
			SELECT novela_id, category_id, COUNT(*) as count
			FROM user_novela_bookmarks
			WHERE novela_id = ANY($1)
			GROUP BY novela_id, category_id
		) s
		JOIN bookmark_categories c ON s.category_id = c.id
		GROUP BY s.novela_id`

	rows, err := tx.QueryContext(ctx, query, pq.Array(novelaIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var novelaID int
		var summary db.NovelaBookmarkSummary
		var distJSON []byte
		var totalCount int

		if err := rows.Scan(&novelaID, &totalCount, &distJSON); err != nil {
			return nil, err
		}

		summary.TotalCount = totalCount
		summary.Distribution = make(map[string]int)
		if len(distJSON) > 0 {
			json.Unmarshal(distJSON, &summary.Distribution)
		}

		result[novelaID] = &summary
	}

	return result, rows.Err()
}

func (r *NovelaBookmarkRepository) GetCategoryByName(ctx context.Context, userID int, name string) (int, error) {
	query := `SELECT id FROM bookmark_categories WHERE name = $1 AND (user_id IS NULL OR user_id = $2) ORDER BY user_id DESC NULLS LAST LIMIT 1`
	var id int
	err := r.DB.QueryRowContext(ctx, query, name, userID).Scan(&id)
	return id, err
}

func (r *NovelaBookmarkRepository) GetCategories(ctx context.Context, userID int) ([]db.BookmarkCategory, error) {
	query := `SELECT id, user_id, name, visible FROM bookmark_categories WHERE user_id IS NULL OR user_id = $1 ORDER BY id ASC`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []db.BookmarkCategory
	for rows.Next() {
		var cat db.BookmarkCategory
		if err := rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Visible); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (r *NovelaBookmarkRepository) GetCategoriesWithCount(ctx context.Context, userID int) ([]db.BookmarkCategoryCount, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.visible, COUNT(b.novela_id) as count
		FROM bookmark_categories c
		LEFT JOIN user_novela_bookmarks b ON c.id = b.category_id AND b.user_id = $1
		WHERE c.user_id IS NULL OR c.user_id = $1
		GROUP BY c.id, c.user_id, c.name, c.visible
		ORDER BY c.id ASC`
	
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []db.BookmarkCategoryCount
	for rows.Next() {
		var cat db.BookmarkCategoryCount
		if err := rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Visible, &cat.Count); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (r *NovelaBookmarkRepository) CreateCategory(ctx context.Context, userID int, name string) (int, error) {
	query := `INSERT INTO bookmark_categories (user_id, name) VALUES ($1, $2) RETURNING id`
	var id int
	err := r.DB.QueryRowContext(ctx, query, userID, name).Scan(&id)
	return id, err
}

func (r *NovelaBookmarkRepository) UpdateCategory(ctx context.Context, id int, userID int, name string, visible bool) error {
	query := `UPDATE bookmark_categories SET name = $1, visible = $2 WHERE id = $3 AND user_id = $4`
	res, err := r.DB.ExecContext(ctx, query, name, visible, id, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category not found or not owned by user")
	}
	return nil
}

func (r *NovelaBookmarkRepository) DeleteCategory(ctx context.Context, id int, userID int) error {
	query := `DELETE FROM bookmark_categories WHERE id = $1 AND user_id = $2`
	res, err := r.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category not found or not owned by user")
	}
	return nil
}

func (r *NovelaBookmarkRepository) GetUsersWithNovelaBookmark(ctx context.Context, novelaID int) ([]int, error) {
	query := `SELECT user_id FROM user_novela_bookmarks WHERE novela_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, novelaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}
