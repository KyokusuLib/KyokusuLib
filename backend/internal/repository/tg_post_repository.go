package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lib/pq"
)

type TgPostRepository struct {
	DB *sql.DB
}

func NewTgPostRepository(db *sql.DB) *TgPostRepository {
	return &TgPostRepository{DB: db}
}

func (r *TgPostRepository) Create(ctx context.Context, p *db.TgPost) (*db.TgPost, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO tg_posts (message_id, text)
		VALUES ($1, $2)
		ON CONFLICT (message_id) DO NOTHING
		RETURNING id, created_at`
	err = tx.QueryRowContext(ctx, query, p.MessageID, p.Text).
		Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(p.Images) > 0 {
		imgQuery := `INSERT INTO tg_post_images (post_id, position, image_path) VALUES ($1, $2, $3)`
		stmt, err := tx.PrepareContext(ctx, imgQuery)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		for _, img := range p.Images {
			if _, err := stmt.ExecContext(ctx, p.ID, img.Position, img.ImagePath); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *TgPostRepository) List(ctx context.Context, limit, offset int) ([]*db.TgPost, error) {
	query := `
		SELECT id, message_id, text, created_at
		FROM tg_posts
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2`
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*db.TgPost
	for rows.Next() {
		var p db.TgPost
		if err := rows.Scan(&p.ID, &p.MessageID, &p.Text, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	if posts == nil {
		return []*db.TgPost{}, nil
	}

	if err := r.attachImages(ctx, posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *TgPostRepository) GetByID(ctx context.Context, id int64) (*db.TgPost, error) {
	query := `
		SELECT id, message_id, text, created_at
		FROM tg_posts
		WHERE id = $1`
	var p db.TgPost
	err := r.DB.QueryRowContext(ctx, query, id).
		Scan(&p.ID, &p.MessageID, &p.Text, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := r.attachImages(ctx, []*db.TgPost{&p}); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *TgPostRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tg_posts WHERE id = $1 RETURNING id`
	var deletedID int64
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&deletedID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TgPostRepository) GetSinceID(ctx context.Context, lastID int64, limit int) ([]*db.TgPost, error) {
	query := `
		SELECT id, message_id, text, created_at
		FROM tg_posts
		WHERE id > $1
		ORDER BY id ASC
		LIMIT $2`
	rows, err := r.DB.QueryContext(ctx, query, lastID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*db.TgPost
	for rows.Next() {
		var p db.TgPost
		if err := rows.Scan(&p.ID, &p.MessageID, &p.Text, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	if posts == nil {
		return []*db.TgPost{}, nil
	}

	if err := r.attachImages(ctx, posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *TgPostRepository) attachImages(ctx context.Context, posts []*db.TgPost) error {
	postIDs := make([]int64, 0, len(posts))
	for _, p := range posts {
		postIDs = append(postIDs, p.ID)
	}

	query := `
		SELECT id, post_id, position, image_path
		FROM tg_post_images
		WHERE post_id = ANY($1)
		ORDER BY post_id, position ASC`
	rows, err := r.DB.QueryContext(ctx, query, pq.Array(postIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	imagesByPost := make(map[int64][]db.TgPostImage)
	for rows.Next() {
		var img db.TgPostImage
		if err := rows.Scan(&img.ID, &img.PostID, &img.Position, &img.ImagePath); err != nil {
			return err
		}
		imagesByPost[img.PostID] = append(imagesByPost[img.PostID], img)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, p := range posts {
		if imgs, ok := imagesByPost[p.ID]; ok {
			p.Images = imgs
		}
	}
	return nil
}
