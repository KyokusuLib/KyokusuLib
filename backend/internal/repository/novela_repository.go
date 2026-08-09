package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lib/pq"
)

type NovelaRepository struct {
	DB *sql.DB
}

func NewNovelaRepository(db *sql.DB) *NovelaRepository {
	return &NovelaRepository{DB: db}
}

func (r *NovelaRepository) Create(ctx context.Context, n *db.Novela) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO novela (
			title, alternative_titles, description, type, age_rating, 
			release_date, status, country, translation_status, poster_url, views
		)
		SELECT 
            $1::text, $2::text[], $3::text, $4::text, $5::text, 
            $6::timestamp, $7::text, $8::text, $9::text, $10::text, $11::int
		WHERE NOT EXISTS (
			SELECT 1 FROM novela 
			WHERE title ILIKE $1::text 
			   OR $1::text ILIKE ANY(alternative_titles)
			   OR title = ANY($2::text[])
			   OR alternative_titles && $2::text[]
		)
		RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		n.Title,
		pq.Array(n.AlternativeTitles),
		n.Description,
		n.Type,
		n.AgeRating,
		n.ReleaseDate,
		n.Status,
		n.Country,
		n.TranslationStatus,
		n.PosterURL,
		0,
	).Scan(&n.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("novela with this title already exists")
		}
		return fmt.Errorf("failed to insert novela: %w", err)
	}

	if len(n.Genres) > 0 {
		if err := r.linkTags(ctx, tx, n.ID, n.Genres, "genres", "novela_genres", "genre_id"); err != nil {
			return err
		}
	}

	if len(n.Categories) > 0 {
		if err := r.linkTags(ctx, tx, n.ID, n.Categories, "categories", "novela_categories", "category_id"); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *NovelaRepository) GetFullByID(tx *sql.Tx, ctx context.Context, id, userID int) (*db.Novela, error) {
	n := &db.Novela{}
	var (
		authorsJSON     []byte
		volumesJSON     []byte
		userRating      sql.NullInt32
		bookmark        sql.NullString
		lastReadChapter sql.NullString
	)

	query := `
		SELECT 
			n.id, 
			n.title, 
			COALESCE(n.alternative_titles, '{}')::text[],
			n.description, 
			n.type, 
			n.age_rating, 
			n.release_date, 
			n.status,
			n.country, 
			n.translation_status, 
			n.poster_url, 
			n.views,
			
			COALESCE((
				SELECT array_agg(g.name) 
				FROM novela_genres ng 
				JOIN genres g ON ng.genre_id = g.id 
				WHERE ng.novela_id = n.id
			), '{}')::text[],

			COALESCE((
				SELECT array_agg(c.name) 
				FROM novela_categories nc 
				JOIN categories c ON nc.category_id = c.id 
				WHERE nc.novela_id = n.id
			), '{}')::text[],

			COALESCE((
				SELECT json_agg(json_build_object('id', a.id, 'name', a.name, 'role', na.role))
				FROM novela_authors na 
				JOIN authors a ON na.author_id = a.id 
				WHERE na.novela_id = n.id
			), '[]'),

			COALESCE((
				SELECT json_agg(
					json_build_object(
						'id', v.id,
						'title', v.title,
						'number', v.volume_number,
						'chapters', COALESCE((
							SELECT json_agg(
								json_build_object(
									'id', ch.id,
									'title', ch.title,
									'number', ch.chapter_number,
									'is_read', CASE WHEN $2 > 0 THEN COALESCE((SELECT true FROM read_chapters rc WHERE rc.chapter_id = ch.id AND rc.user_id = $2), false) ELSE false END,
									'images', COALESCE((
										SELECT json_agg(
											json_build_object('id', img.id, 'image_url', img.image_url, 'caption', img.caption, 'position', img.position)
											ORDER BY img.position ASC
										) FROM novela_chapter_images img WHERE img.chapter_id = ch.id
									), '[]'::json)
								) ORDER BY ch.chapter_number
							) FROM novela_chapters ch WHERE ch.novela_volume_id = v.id AND ch.status = 'approved'
						), '[]'::json)
					) ORDER BY v.volume_number
				) FROM novela_volumes v WHERE v.novela_id = n.id AND v.status = 'approved'
			), '[]'),

			CASE 
				WHEN $2 > 0 THEN (SELECT c.name FROM user_novela_bookmarks b JOIN bookmark_categories c ON b.category_id = c.id WHERE b.novela_id = n.id AND b.user_id = $2)
				ELSE NULL 
			END as user_category,

			CASE 
				WHEN $2 > 0 THEN COALESCE((SELECT has_liked FROM user_novela_likes WHERE novela_id = n.id AND user_id = $2), FALSE)
				ELSE FALSE 
			END as has_liked,

			COALESCE((SELECT COUNT(*) FROM user_novela_likes WHERE novela_id = n.id AND has_liked = TRUE), 0),

			CASE
				WHEN $2 > 0 THEN (SELECT rating FROM novela_ratings WHERE novela_id = n.id AND user_id = $2)
				ELSE 0
			END as user_rating,

			CASE
				WHEN $2 > 0 THEN (
					SELECT json_build_object('id', c.id, 'number', c.chapter_number, 'scroll_position', rc.scroll_position)::text
					FROM read_chapters rc
					JOIN novela_chapters c ON rc.chapter_id = c.id
					JOIN novela_volumes v ON c.novela_volume_id = v.id
					WHERE v.novela_id = n.id AND rc.user_id = $2
					ORDER BY rc.created_at DESC
					LIMIT 1
				)
				ELSE NULL
			END as last_read_chapter

		FROM novela n
		WHERE n.id = $1`

	err := tx.QueryRowContext(ctx, query, id, userID).Scan(
		&n.ID,
		&n.Title,
		pq.Array(&n.AlternativeTitles),
		&n.Description,
		&n.Type,
		&n.AgeRating,
		&n.ReleaseDate,
		&n.Status,
		&n.Country,
		&n.TranslationStatus,
		&n.PosterURL,
		&n.Views,
		pq.Array(&n.Genres),
		pq.Array(&n.Categories),
		&authorsJSON,
		&volumesJSON,
		&bookmark,
		&n.HasLiked,
		&n.LikeCount,
		&userRating,
		&lastReadChapter,
	)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(authorsJSON) > 0 {
		if err := json.Unmarshal(authorsJSON, &n.Authors); err != nil {
			return nil, err
		}
	}
	if len(volumesJSON) > 0 {
		if err := json.Unmarshal(volumesJSON, &n.Volumes); err != nil {
			return nil, err
		}
	}

	if bookmark.Valid {
		n.Bookmark = &db.BookmarkCategory{Name: bookmark.String}
	}

	if userRating.Valid {
		n.UserRating = int(userRating.Int32)
	}

	if lastReadChapter.Valid {
		var lc db.LastReadChapter
		if err := json.Unmarshal([]byte(lastReadChapter.String), &lc); err != nil {
			return nil, err
		}
		n.LastReadChapter = &lc
	}

	return n, nil
}

func (r *NovelaRepository) UpdateFull(ctx context.Context, n *db.Novela) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE novela 
		SET title = $1, alternative_titles = $2, description = $3, type = $4, 
		    age_rating = $5, release_date = $6, status = $7, country = $8, 
		    translation_status = $9, poster_url = COALESCE(NULLIF($10, ''), poster_url),
			updated_at = NOW()
		WHERE id = $11`

	res, err := tx.ExecContext(ctx, query,
		n.Title, pq.Array(n.AlternativeTitles), n.Description, n.Type,
		n.AgeRating, n.ReleaseDate, n.Status, n.Country,
		n.TranslationStatus, n.PosterURL, n.ID,
	)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	for _, table := range []string{"novela_genres", "novela_categories", "novela_authors"} {
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE novela_id = $1", table), n.ID); err != nil {
			return err
		}
	}

	if err := r.linkTags(ctx, tx, n.ID, n.Genres, "genres", "novela_genres", "genre_id"); err != nil {
		return err
	}

	if err := r.linkTags(ctx, tx, n.ID, n.Categories, "categories", "novela_categories", "category_id"); err != nil {
		return err
	}

	for _, author := range n.Authors {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO novela_authors (novela_id, author_id, role) VALUES ($1, $2, $3)`,
			n.ID, author.ID, "Author",
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *NovelaRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM novela WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

func (r *NovelaRepository) LinkAuthor(ctx context.Context, novelaID int, authorName string, role string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var authorID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO authors (name) VALUES ($1) 
		ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name 
		RETURNING id`, authorName).Scan(&authorID)

	if err != nil {
		return fmt.Errorf("failed to upsert author: %w", err)
	}

	query := `
		INSERT INTO novela_authors (novela_id, author_id, role) 
		VALUES ($1, $2, $3)
		ON CONFLICT (novela_id, author_id) 
		DO UPDATE SET role = EXCLUDED.role`

	_, err = tx.ExecContext(ctx, query, novelaID, authorID, role)
	if err != nil {
		return fmt.Errorf("failed to link author: %w", err)
	}

	return tx.Commit()
}

func (r *NovelaRepository) linkTags(ctx context.Context, tx *sql.Tx, novelaID int, tags []string, tableName, linkTable, fkCol string) error {
	for _, name := range tags {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		var tagID int
		upsertQuery := fmt.Sprintf(`
			INSERT INTO %s (name) VALUES (TRIM($1)) 
			ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name 
			RETURNING id`, tableName)

		err := tx.QueryRowContext(ctx, upsertQuery, name).Scan(&tagID)
		if err != nil {
			return err
		}

		linkQuery := fmt.Sprintf(`
			INSERT INTO %s (novela_id, %s) VALUES ($1, $2) 
			ON CONFLICT DO NOTHING`, linkTable, fkCol)

		_, err = tx.ExecContext(ctx, linkQuery, novelaID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *NovelaRepository) CreateVolume(ctx context.Context, novelaID int, title string, number int) (string, error) {
	query := `
		INSERT INTO novela_volumes (novela_id, title, volume_number) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	var id string
	err := r.DB.QueryRowContext(ctx, query, novelaID, title, number).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create volume: %w", err)
	}
	return id, nil
}

func (r *NovelaRepository) CreateChapter(ctx context.Context, volumeID string, ch *db.NovelaChapter) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO novela_chapters (novela_volume_id, title, chapter_number, content) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	err = tx.QueryRowContext(ctx, query, volumeID, ch.Title, ch.Number, ch.Content).Scan(&ch.ID)
	if err != nil {
		return fmt.Errorf("failed to create chapter: %w", err)
	}

	if len(ch.Images) > 0 {
		imgQuery := `INSERT INTO novela_chapter_images (chapter_id, image_url, caption, position) VALUES ($1, $2, $3, $4)`
		stmt, err := tx.PrepareContext(ctx, imgQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, img := range ch.Images {
			if _, err := stmt.ExecContext(ctx, ch.ID, img.ImageURL, img.Caption, img.Position); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *NovelaRepository) GetChapterByID(ctx context.Context, chapterID string) (*db.NovelaChapter, error) {
	ch := &db.NovelaChapter{}

	query := `SELECT id, novela_volume_id, title, chapter_number, content FROM novela_chapters WHERE id = $1`
	err := r.DB.QueryRowContext(ctx, query, chapterID).Scan(&ch.ID, &ch.VolumeID, &ch.Title, &ch.Number, &ch.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	imgQuery := `SELECT id, image_url, caption, position FROM novela_chapter_images WHERE chapter_id = $1 ORDER BY position ASC`
	rows, err := r.DB.QueryContext(ctx, imgQuery, chapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img db.NovelaChapterImage
		var caption sql.NullString

		if err := rows.Scan(&img.ID, &img.ImageURL, &caption, &img.Position); err != nil {
			return nil, err
		}
		img.Caption = caption.String
		ch.Images = append(ch.Images, img)
	}

	return ch, nil
}

func (r *NovelaRepository) GetChapterReaderDetails(ctx context.Context, chapterID string, userID int) (*dto.ChapterReaderResponse, error) {
	res := &dto.ChapterReaderResponse{}

	query := `
		WITH current_chapter AS (
			SELECT c.id, c.chapter_number, v.volume_number, v.novela_id, n.title as novela_title
			FROM novela_chapters c
			JOIN novela_volumes v ON c.novela_volume_id = v.id
			JOIN novela n ON v.novela_id = n.id
			WHERE c.id = $1
		)
		SELECT 
			c.id, c.title, c.chapter_number, c.content, c.novela_volume_id, cc.novela_id, cc.novela_title, cc.volume_number,
			COALESCE((SELECT rc.scroll_position FROM read_chapters rc WHERE rc.chapter_id = c.id AND rc.user_id = $2), 0),
			(SELECT c2.id 
			 FROM novela_chapters c2 
			 JOIN novela_volumes v2 ON c2.novela_volume_id = v2.id 
			 CROSS JOIN current_chapter cc
			 WHERE v2.novela_id = cc.novela_id 
			   AND (v2.volume_number < cc.volume_number OR (v2.volume_number = cc.volume_number AND c2.chapter_number < cc.chapter_number))
			   AND c2.status = 'approved' AND v2.status = 'approved'
			 ORDER BY v2.volume_number DESC, c2.chapter_number DESC
			 LIMIT 1) as prev_id,
			 
			(SELECT c3.id 
			 FROM novela_chapters c3 
			 JOIN novela_volumes v3 ON c3.novela_volume_id = v3.id 
			 CROSS JOIN current_chapter cc
			 WHERE v3.novela_id = cc.novela_id 
			   AND (v3.volume_number > cc.volume_number OR (v3.volume_number = cc.volume_number AND c3.chapter_number > cc.chapter_number))
			   AND c3.status = 'approved' AND v3.status = 'approved'
			 ORDER BY v3.volume_number ASC, c3.chapter_number ASC
			 LIMIT 1) as next_id
		FROM novela_chapters c
		JOIN current_chapter cc ON c.id = cc.id`

	err := r.DB.QueryRowContext(ctx, query, chapterID, userID).Scan(
		&res.ID, &res.Title, &res.Number, &res.Content, &res.VolumeID, &res.NovelaID, &res.NovelaTitle, &res.VolumeNumber,
		&res.ScrollPosition,
		&res.PrevChapterID, &res.NextChapterID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	imgQuery := `SELECT id, image_url, caption, position FROM novela_chapter_images WHERE chapter_id = $1 ORDER BY position ASC`
	rows, err := r.DB.QueryContext(ctx, imgQuery, chapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img dto.NovelaChapterImage
		var caption sql.NullString

		if err := rows.Scan(&img.ID, &img.ImageURL, &caption, &img.Position); err != nil {
			return nil, err
		}
		img.Caption = caption.String
		res.Images = append(res.Images, img)
	}

	return res, nil
}

func (r *NovelaRepository) UpdatePoster(ctx context.Context, id int, posterURL string) error {
	query := `UPDATE novela SET poster_url = $2 WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id, posterURL)
	return err
}

func (r *NovelaRepository) GetIDByTitle(ctx context.Context, title string) (int, error) {
	query := `
		SELECT id FROM novela 
		WHERE title ILIKE $1 
		   OR $1 ILIKE ANY(alternative_titles)
		LIMIT 1`
	var id int
	err := r.DB.QueryRowContext(ctx, query, title).Scan(&id)
	return id, err
}

func (r *NovelaRepository) GetNovelas(tx *sql.Tx, ctx context.Context, userID int, f dto.NovelaFilters) ([]db.Novela, int, error) {
	var args []interface{}
	where := []string{"$1::int IS NOT NULL"}
	argID := 1

	args = append(args, userID)
	argID++

	if f.Search != "" {
		where = append(where, fmt.Sprintf("(n.title ILIKE $%d::text OR $%d::text ILIKE ANY(n.alternative_titles))", argID, argID))
		args = append(args, "%"+f.Search+"%")
		argID++
	}

	if len(f.Genres) > 0 {
		where = append(where, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM novela_genres ng JOIN genres g ON ng.genre_id = g.id 
			WHERE ng.novela_id = n.id AND g.name = ANY($%d)
		)`, argID))
		args = append(args, pq.Array(f.Genres))
		argID++
	}

	if len(f.Categories) > 0 {
		where = append(where, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM novela_categories nc JOIN categories c ON nc.category_id = c.id 
			WHERE nc.novela_id = n.id AND c.name = ANY($%d)
		)`, argID))
		args = append(args, pq.Array(f.Categories))
		argID++
	}

	if len(f.Country) > 0 {
		where = append(where, fmt.Sprintf(`n.country = ANY($%d)`, argID))
		args = append(args, pq.Array(f.Country))
		argID++
	}

	if f.Type != "" {
		where = append(where, fmt.Sprintf("n.type = $%d", argID))
		args = append(args, f.Type)
		argID++
	}

	if f.Status != "" {
		where = append(where, fmt.Sprintf("n.status = $%d", argID))
		args = append(args, f.Status)
		argID++
	}

	if f.TranslationStatus != "" {
		where = append(where, fmt.Sprintf("n.translation_status = $%d", argID))
		args = append(args, f.TranslationStatus)
		argID++
	}

	if f.AgeRating != "" {
		where = append(where, fmt.Sprintf("n.age_rating = $%d", argID))
		args = append(args, f.AgeRating)
		argID++
	}

	if f.ChaptersFrom > 0 {
		where = append(where, fmt.Sprintf(`(
			SELECT COUNT(*) FROM novela_chapters ch 
			JOIN novela_volumes v ON ch.novela_volume_id = v.id 
			WHERE v.novela_id = n.id AND ch.status = 'approved' AND v.status = 'approved'
		) >= $%d`, argID))
		args = append(args, f.ChaptersFrom)
		argID++
	}

	if f.ChaptersTo > 0 {
		where = append(where, fmt.Sprintf(`(
			SELECT COUNT(*) FROM novela_chapters ch 
			JOIN novela_volumes v ON ch.novela_volume_id = v.id 
			WHERE v.novela_id = n.id AND ch.status = 'approved' AND v.status = 'approved'
		) <= $%d`, argID))
		args = append(args, f.ChaptersTo)
		argID++
	}

	if f.YearFrom > 0 {
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM n.release_date) >= $%d", argID))
		args = append(args, f.YearFrom)
		argID++
	}

	if f.YearTo > 0 {
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM n.release_date) <= $%d", argID))
		args = append(args, f.YearTo)
		argID++
	}

	if f.AuthorID > 0 {
		where = append(where, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM novela_authors na 
			WHERE na.novela_id = n.id AND na.author_id = $%d
		)`, argID))
		args = append(args, f.AuthorID)
		argID++
	}

	orderBy := "n.created_at DESC"
	switch f.Sort {
	case dto.SortPopular:
		orderBy = "avg_rating DESC NULLS LAST"
	case dto.SortTrending:
		orderBy = "n.views DESC"
	case dto.SortUpdated:
		orderBy = "last_chapter_at DESC NULLS LAST"
	case dto.SortNew:
		orderBy = "n.created_at DESC"
	}

	query := fmt.Sprintf(`
    SELECT 
        n.id, n.title, n.description, n.age_rating, n.release_date, n.translation_status, n.country, COALESCE(n.alternative_titles, '{}')::text[], n.type, n.status, n.poster_url, n.views,
        COALESCE((SELECT array_agg(g.name) FROM novela_genres ng JOIN genres g ON ng.genre_id = g.id WHERE ng.novela_id = n.id), '{}')::text[],
        COALESCE((SELECT array_agg(c.name) FROM novela_categories nc JOIN categories c ON nc.category_id = c.id WHERE nc.novela_id = n.id), '{}')::text[],
        COALESCE((SELECT AVG(rating) FROM novela_ratings WHERE novela_id = n.id), 0) as avg_rating,
        (SELECT MAX(ch.created_at) FROM novela_chapters ch JOIN novela_volumes v ON ch.novela_volume_id = v.id WHERE v.novela_id = n.id AND ch.status = 'approved' AND v.status = 'approved') as last_chapter_at,
        
        (SELECT c.name FROM user_novela_bookmarks b JOIN bookmark_categories c ON b.category_id = c.id WHERE b.novela_id = n.id AND b.user_id = $1) as user_bookmark,
        
        COALESCE((SELECT has_liked FROM user_novela_likes WHERE novela_id = n.id AND user_id = $1), FALSE) as has_liked,
        COALESCE((SELECT COUNT(*) FROM user_novela_likes WHERE novela_id = n.id), 0) as like_count
    FROM novela n
    WHERE %s
    ORDER BY %s
    LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "), orderBy, argID, argID+1)

	finalArgs := append(args, f.Limit, f.Offset)
	rows, err := tx.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var novelas []db.Novela
	for rows.Next() {
		var n db.Novela
		var lastChapterAt sql.NullTime
		var bookmark sql.NullString
		err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Description,
			&n.AgeRating,
			&n.ReleaseDate,
			&n.TranslationStatus,
			&n.Country,
			pq.Array(&n.AlternativeTitles),
			&n.Type,
			&n.Status,
			&n.PosterURL,
			&n.Views,
			pq.Array(&n.Genres),
			pq.Array(&n.Categories),
			&n.Rating,
			&lastChapterAt,
			&bookmark,
			&n.HasLiked,
			&n.LikeCount,
		)
		if err != nil {
			return nil, 0, err
		}

		if bookmark.Valid {
			n.Bookmark = &db.BookmarkCategory{Name: bookmark.String}
		}

		novelas = append(novelas, n)
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM novela n WHERE %s", strings.Join(where, " AND "))
	if err := tx.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return novelas, total, nil
}

func (r *NovelaRepository) GetUserNovelaBookmarks(ctx context.Context, userID int, categoryID int) ([]db.UserNovelaBookmark, error) {
	query := `
		SELECT n.id, n.title, n.poster_url, n.type,
		       COALESCE((SELECT AVG(rating) FROM novela_ratings WHERE novela_id = n.id), 0) as rating
		FROM novela n
		JOIN user_novela_bookmarks b ON n.id = b.novela_id
		WHERE b.user_id = $1 AND b.category_id = $2
		ORDER BY b.updated_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var novelas []db.UserNovelaBookmark
	for rows.Next() {
		var n db.UserNovelaBookmark
		if err := rows.Scan(&n.ID, &n.Title, &n.PosterURL, &n.Type, &n.Rating); err != nil {
			return nil, err
		}
		novelas = append(novelas, n)
	}

	return novelas, nil
}

func (r *NovelaRepository) CountNovelaTeams(ctx context.Context, novelaID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM novela_teams WHERE novela_id = $1`
	err := r.DB.QueryRowContext(ctx, query, novelaID).Scan(&count)
	return count, err
}

func (r *NovelaRepository) AddTeamToNovela(ctx context.Context, novelaID, teamID int) error {
	query := `INSERT INTO novela_teams (novela_id, team_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.DB.ExecContext(ctx, query, novelaID, teamID)
	return err
}

func (r *NovelaRepository) AddVolume(ctx context.Context, novelaID int, volumeNumber int, title string, status string, userID int) (string, error) {
	query := `INSERT INTO novela_volumes (novela_id, volume_number, title, status, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := r.DB.QueryRowContext(ctx, query, novelaID, volumeNumber, title, status, userID).Scan(&id)
	return id, err
}

func (r *NovelaRepository) AddChapter(ctx context.Context, volumeID string, chapterNumber float64, title, content string, status string, userID int) (string, error) {
	query := `INSERT INTO novela_chapters (novela_volume_id, chapter_number, title, content, status, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id string
	err := r.DB.QueryRowContext(ctx, query, volumeID, chapterNumber, title, content, status, userID).Scan(&id)
	return id, err
}

func (r *NovelaRepository) DeleteChapterImages(ctx context.Context, chapterID string) error {
	query := `DELETE FROM novela_chapter_images WHERE chapter_id = $1`
	_, err := r.DB.ExecContext(ctx, query, chapterID)
	return err
}

func (r *NovelaRepository) AddChapterImage(ctx context.Context, chapterID string, imageURL, caption string, position int) (int, error) {
	query := `INSERT INTO novela_chapter_images (chapter_id, image_url, caption, position) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.DB.QueryRowContext(ctx, query, chapterID, imageURL, caption, position).Scan(&id)
	return id, err
}

func (r *NovelaRepository) UpdateChapter(ctx context.Context, chapterID string, chapterNumber float64, title, content string) error {
	query := `UPDATE novela_chapters SET chapter_number = $1, title = $2, content = $3 WHERE id = $4`
	_, err := r.DB.ExecContext(ctx, query, chapterNumber, title, content, chapterID)
	return err
}

func (r *NovelaRepository) DeleteChapter(ctx context.Context, chapterID string) error {
	query := `DELETE FROM novela_chapters WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, chapterID)
	return err
}

func (r *NovelaRepository) DeleteVolume(ctx context.Context, volumeID string) error {
	query := `DELETE FROM novela_volumes WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, volumeID)
	return err
}

func (r *NovelaRepository) GetTitleByID(ctx context.Context, id int) (string, error) {
	query := `SELECT title FROM novela WHERE id = $1`
	var title string
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&title)
	return title, err
}

func (r *NovelaRepository) GetVolumeIDByNumber(ctx context.Context, novelaID int, volumeNumber int) (string, error) {
	query := `SELECT id FROM novela_volumes WHERE novela_id = $1 AND volume_number = $2`
	var id string
	err := r.DB.QueryRowContext(ctx, query, novelaID, volumeNumber).Scan(&id)
	return id, err
}

func (r *NovelaRepository) GetChapterIDByNumber(ctx context.Context, volumeID string, chapterNumber float64) (string, error) {
	query := `SELECT id FROM novela_chapters WHERE novela_volume_id = $1 AND chapter_number = $2`
	var id string
	err := r.DB.QueryRowContext(ctx, query, volumeID, chapterNumber).Scan(&id)
	return id, err
}

func (r *NovelaRepository) GetNovelaIDByVolumeID(ctx context.Context, volumeID string) (int, error) {
	query := `SELECT novela_id FROM novela_volumes WHERE id = $1`
	var novelaID int
	err := r.DB.QueryRowContext(ctx, query, volumeID).Scan(&novelaID)
	return novelaID, err
}

func (r *NovelaRepository) CheckUserNovelaTeamPermission(ctx context.Context, userID, novelaID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM novela_teams nt
			JOIN team_members tm ON nt.team_id = tm.team_id
			WHERE nt.novela_id = $1 
			  AND tm.user_id = $2 
			  AND tm.role IN ('owner', 'admin', 'moderator')
		)`
	var exists bool
	err := r.DB.QueryRowContext(ctx, query, novelaID, userID).Scan(&exists)
	return exists, err
}

func (r *NovelaRepository) GetPendingVolumes(ctx context.Context) ([]dto.PendingVolume, error) {
	query := `
		SELECT v.id, v.novela_id, n.title, v.volume_number, COALESCE(v.title, ''),
		       v.status, COALESCE(v.created_by, 0), COALESCE(up.name, '')
		FROM novela_volumes v
		JOIN novela n ON n.id = v.novela_id
		LEFT JOIN user_profiles up ON up.user_id = v.created_by
		WHERE v.status = 'pending'
		ORDER BY v.volume_number`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	volumes := make([]dto.PendingVolume, 0)
	for rows.Next() {
		var v dto.PendingVolume
		if err := rows.Scan(
			&v.ID,
			&v.NovelaID,
			&v.NovelaTitle,
			&v.Number,
			&v.Title,
			&v.Status,
			&v.CreatedBy,
			&v.CreatedByName,
		); err != nil {
			return nil, err
		}
		volumes = append(volumes, v)
	}
	return volumes, nil
}

func (r *NovelaRepository) GetPendingChapters(ctx context.Context) ([]dto.PendingChapter, error) {
	query := `
		SELECT c.id, c.novela_volume_id, v.volume_number, v.novela_id, n.title,
		       c.chapter_number, COALESCE(c.title, ''), COALESCE(c.content, ''),
		       c.status, COALESCE(c.created_by, 0), COALESCE(up.name, '')
		FROM novela_chapters c
		JOIN novela_volumes v ON v.id = c.novela_volume_id
		JOIN novela n ON n.id = v.novela_id
		LEFT JOIN user_profiles up ON up.user_id = c.created_by
		WHERE c.status = 'pending'
		ORDER BY v.volume_number, c.chapter_number`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chapters := make([]dto.PendingChapter, 0)
	for rows.Next() {
		var ch dto.PendingChapter
		if err := rows.Scan(
			&ch.ID,
			&ch.VolumeID,
			&ch.VolumeNumber,
			&ch.NovelaID,
			&ch.NovelaTitle,
			&ch.Number,
			&ch.Title,
			&ch.Content,
			&ch.Status,
			&ch.CreatedBy,
			&ch.CreatedByName,
		); err != nil {
			return nil, err
		}
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

func (r *NovelaRepository) UpdateVolume(ctx context.Context, id string, volumeNumber int, title string) error {
	query := `UPDATE novela_volumes SET volume_number = $1, title = $2 WHERE id = $3`
	_, err := r.DB.ExecContext(ctx, query, volumeNumber, title, id)
	return err
}

func (r *NovelaRepository) UpdateVolumeStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE novela_volumes SET status = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, id)
	return err
}

func (r *NovelaRepository) UpdateChapterStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE novela_chapters SET status = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, id)
	return err
}

func (r *NovelaRepository) GetVolumeByID(ctx context.Context, id string) (*db.NovelaVolume, error) {
	v := &db.NovelaVolume{}
	var novelaID int
	query := `SELECT id, novela_id, volume_number, title, status, created_by FROM novela_volumes WHERE id = $1`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&v.ID, &novelaID, &v.Number, &v.Title, &v.Status, &v.CreatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (r *NovelaRepository) GetNovelaIDByChapterID(ctx context.Context, chapterID string) (int, error) {
	query := `
		SELECT v.novela_id 
		FROM novela_chapters c
		JOIN novela_volumes v ON c.novela_volume_id = v.id
		WHERE c.id = $1`
	var novelaID int
	err := r.DB.QueryRowContext(ctx, query, chapterID).Scan(&novelaID)
	return novelaID, err
}

func (r *NovelaRepository) GetVolumeIDByChapterID(ctx context.Context, chapterID string) (string, error) {
	query := `SELECT novela_volume_id FROM novela_chapters WHERE id = $1`
	var volumeID string
	err := r.DB.QueryRowContext(ctx, query, chapterID).Scan(&volumeID)
	return volumeID, err
}

func (r *NovelaRepository) MarkChapterAsRead(ctx context.Context, userID int, chapterID string) (bool, error) {
	query := `
		INSERT INTO read_chapters (user_id, chapter_id) 
		VALUES ($1, $2) 
		ON CONFLICT (user_id, chapter_id) DO NOTHING`
	result, err := r.DB.ExecContext(ctx, query, userID, chapterID)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *NovelaRepository) GetUserReadChaptersByNovela(ctx context.Context, userID, novelaID int) ([]db.ReadChapter, error) {
	query := `
		SELECT rc.user_id, rc.chapter_id, rc.scroll_position, rc.created_at
		FROM read_chapters rc
		JOIN novela_chapters c ON rc.chapter_id = c.id
		JOIN novela_volumes v ON c.novela_volume_id = v.id
		WHERE v.novela_id = $1 AND rc.user_id = $2
		ORDER BY rc.created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, novelaID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []db.ReadChapter
	for rows.Next() {
		var ch db.ReadChapter
		if err := rows.Scan(&ch.UserID, &ch.ChapterID, &ch.ScrollPosition, &ch.CreatedAt); err != nil {
			return nil, err
		}
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

func (r *NovelaRepository) UpdateChapterReadPosition(ctx context.Context, userID int, chapterID string, scrollPosition int) error {
	query := `
		UPDATE read_chapters
		SET scroll_position = $3
		WHERE user_id = $1 AND chapter_id = $2`
	_, err := r.DB.ExecContext(ctx, query, userID, chapterID, scrollPosition)
	return err
}

func (r NovelaRepository) IsExistChapter(ctx context.Context, chapterID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM novela_chapters WHERE id = $1)`
	var exists bool
	err := r.DB.QueryRowContext(ctx, query, chapterID).Scan(&exists)
	return exists, err
}

func (r NovelaRepository) IsReadChapter(ctx context.Context, userID int, chapterID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM read_chapters WHERE user_id = $1 AND chapter_id = $2)`
	var exists bool
	err := r.DB.QueryRowContext(ctx, query, userID, chapterID).Scan(&exists)
	return exists, err
}

func (r NovelaRepository) GetMostSearchedGenres(ctx context.Context, limit int) ([]string, error) {
	query := `
		SELECT 
			g.name
		FROM genres g
			JOIN novela_genres ng ON ng.genre_id = g.id
			JOIN novela_volumes v ON ng.novela_id = v.novela_id
			JOIN novela_chapters c ON v.id = c.novela_volume_id
			JOIN read_chapters rc ON c.id = rc.chapter_id
		GROUP BY g.id, g.name
		ORDER BY COUNT(c.id) DESC
		LIMIT $1;
	`
	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]string, 0, limit)

	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	return genres, nil
}

func (r NovelaRepository) GetMostSearchedCategories(ctx context.Context, limit int) ([]string, error) {
	query := `
		SELECT 
			c.name
		FROM categories c
			JOIN novela_categories nc ON nc.category_id = c.id
			JOIN novela_volumes v ON nc.novela_id = v.novela_id
			JOIN novela_chapters ch ON v.id = ch.novela_volume_id
			JOIN read_chapters rc ON ch.id = rc.chapter_id
		GROUP BY c.id, c.name
		ORDER BY COUNT(ch.id) DESC
		LIMIT $1;
	`
	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]string, 0, limit)

	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
