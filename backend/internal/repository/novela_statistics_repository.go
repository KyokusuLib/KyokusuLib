package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lanxre/kyokusulib/internal/constants"
	"github.com/lanxre/kyokusulib/internal/models/db"
)

type NovelaStatisticsRepository struct {
	DB *sql.DB
}

func NewNovelaStatisticsRepository(db *sql.DB) *NovelaStatisticsRepository {
	return &NovelaStatisticsRepository{DB: db}
}

func (r *NovelaStatisticsRepository) GetTotalStatistics(
	ctx context.Context,
	limit int,
	offset int,
	period constants.StatisticsPeriodSort,
) ([]db.TotalNovelaStatistics, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	startTime := r.startOfPeriod(period)

	query := `
		WITH bookmark_stats AS (
			SELECT
				novela_id,
				COUNT(*) AS total_bookmark_count
			FROM user_novela_bookmarks
			WHERE ($2::TIMESTAMPTZ IS NULL OR created_at >= $2)
			GROUP BY novela_id
		),
		rating_stats AS (
			SELECT
				novela_id,
				COUNT(*) AS total_rating_count,
				AVG(rating)::float8 AS rating
			FROM novela_ratings
			WHERE ($2::TIMESTAMPTZ IS NULL OR created_at >= $2)
			GROUP BY novela_id
		),
		comment_stats AS (
			SELECT
				novela_id,
				COUNT(*) AS total_comment_count
			FROM novela_comments
			WHERE ($2::TIMESTAMPTZ IS NULL OR created_at >= $2)
			GROUP BY novela_id
		),
		volume_stats AS (
			SELECT
				novela_id,
				COUNT(*) AS total_volume_count
			FROM novela_volumes
			GROUP BY novela_id
		),
		chapter_stats AS (
			SELECT
				nv.novela_id,
				COUNT(*) AS total_chapter_count
			FROM novela_volumes nv
			JOIN novela_chapters nch
				ON nch.novela_volume_id = nv.id
			GROUP BY nv.novela_id
		),
		read_stats AS (
			SELECT
				nv.novela_id,
				COUNT(rc.user_id) AS total_read_count
			FROM read_chapters rc
			JOIN novela_chapters nch
				ON nch.id = rc.chapter_id
			JOIN novela_volumes nv
				ON nv.id = nch.novela_volume_id
			WHERE ($2::TIMESTAMPTZ IS NULL OR rc.created_at >= $2)
			GROUP BY nv.novela_id
		)
		SELECT
			n.id,
			n.title,
			n.description,
			n.poster_url,
			n.release_date,
		
			COALESCE(bs.total_bookmark_count, 0) AS total_bookmark_count,
			COALESCE(rs_read.total_read_count, 0) AS total_read_count,
			COALESCE(rs_rate.total_rating_count, 0) AS total_rating_count,
			COALESCE(cs.total_comment_count, 0) AS total_comment_count,
			COALESCE(vs.total_volume_count, 0) AS total_volume_count,
			COALESCE(chs.total_chapter_count, 0) AS total_chapter_count,
			COALESCE(rs_rate.rating, 0.0) AS rating
		
		FROM novela n
		
		LEFT JOIN bookmark_stats bs
			ON bs.novela_id = n.id
		
		LEFT JOIN read_stats rs_read
			ON rs_read.novela_id = n.id
		
		LEFT JOIN rating_stats rs_rate
			ON rs_rate.novela_id = n.id
		
		LEFT JOIN comment_stats cs
			ON cs.novela_id = n.id
		
		LEFT JOIN volume_stats vs
			ON vs.novela_id = n.id
		
		LEFT JOIN chapter_stats chs
			ON chs.novela_id = n.id
		
		ORDER BY
			rating DESC,
			total_read_count DESC,
			total_bookmark_count DESC,
			total_comment_count DESC,
			total_volume_count DESC,
			total_chapter_count DESC,
			n.id DESC
		LIMIT $1
		OFFSET $3
		`

	rows, err := r.DB.QueryContext(ctx, query, limit, startTime, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]db.TotalNovelaStatistics, 0, limit)

	for rows.Next() {
		var s db.TotalNovelaStatistics
		var releaseDate sql.NullTime

		if err := rows.Scan(
			&s.Novela.NovelaID,
			&s.Novela.Title,
			&s.Novela.Description,
			&s.Novela.PosterURL,
			&releaseDate,
			&s.TotalBookmarkCount,
			&s.TotalReadCount,
			&s.TotalRatingCount,
			&s.TotalCommentCount,
			&s.TotalVolumeCount,
			&s.TotalChapterCount,
			&s.Rating,
		); err != nil {
			return nil, err
		}

		if releaseDate.Valid {
			s.Novela.ReleaseDate = releaseDate.Time
		}

		stats = append(stats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *NovelaStatisticsRepository) GeneralStatistics(
	ctx context.Context, 
	period constants.StatisticsPeriodSort,
) (db.GeneralStatistics, error) {
	
	var stats db.GeneralStatistics

	startTime := r.startOfPeriod(period)

	query := `
		WITH bookmark_stats AS (
			SELECT COUNT(*) AS total
			FROM user_novela_bookmarks
			WHERE ($1::TIMESTAMPTZ IS NULL OR created_at >= $1)
		),
		read_stats AS (
			SELECT COUNT(rc.user_id) AS total
			FROM read_chapters rc
			WHERE ($1::TIMESTAMPTZ IS NULL OR rc.created_at >= $1)
		),
		rating_stats AS (
			SELECT COUNT(*) AS total
			FROM novela_ratings
			WHERE ($1::TIMESTAMPTZ IS NULL OR created_at >= $1)
		),
		comment_stats AS (
			SELECT COUNT(*) AS total
			FROM novela_comments
			WHERE ($1::TIMESTAMPTZ IS NULL OR created_at >= $1)
		),
		volume_stats AS (
			SELECT COUNT(*) AS total
			FROM novela_volumes
		),
		chapter_stats AS (
			SELECT COUNT(*) AS total
			FROM novela_chapters nch
			JOIN novela_volumes nv ON nv.id = nch.novela_volume_id
		)
		SELECT
			COALESCE(bs.total, 0),
			COALESCE(rs.total, 0),
			COALESCE(rts.total, 0),
			COALESCE(cs.total, 0),
			COALESCE(vs.total, 0),
			COALESCE(chs.total, 0)
		FROM bookmark_stats bs
		CROSS JOIN read_stats rs
		CROSS JOIN rating_stats rts
		CROSS JOIN comment_stats cs
		CROSS JOIN volume_stats vs
		CROSS JOIN chapter_stats chs
	`

	if err := r.DB.QueryRowContext(ctx, query, startTime).Scan(
		&stats.GeneralBookmarkCount,
		&stats.GeneralReadCount,
		&stats.GeneralRatingCount,
		&stats.GeneralCommentCount,
		&stats.GeneralVolumeCount,
		&stats.GeneralChapterCount,
	); err != nil {
		return db.GeneralStatistics{}, err
	}

	return stats, nil
}

func (r *NovelaStatisticsRepository) GetMonthlyNovelaSeries(ctx context.Context, limit int) ([]db.NovelaMonthlyRow, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	query := `
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', NOW()),
				date_trunc('year', NOW()) + INTERVAL '11 months',
				'1 month'
			)::timestamp AS month_start
		),
		top_ids AS (
			SELECT v.novela_id, COUNT(*) as total_reads
			FROM read_chapters rc
			JOIN novela_chapters c ON c.id = rc.chapter_id
			JOIN novela_volumes v ON v.id = c.novela_volume_id
			GROUP BY v.novela_id
			ORDER BY total_reads DESC
			LIMIT $1
		),
		monthly_stats AS (
			SELECT 
				v.novela_id, 
				date_trunc('month', rc.created_at) AS month_start, 
				COUNT(*)::int AS read_count
			FROM read_chapters rc
			JOIN novela_chapters c ON c.id = rc.chapter_id
			JOIN novela_volumes v ON v.id = c.novela_volume_id
			WHERE v.novela_id IN (SELECT novela_id FROM top_ids)
			  AND rc.created_at >= date_trunc('year', NOW())
			GROUP BY v.novela_id, month_start
		)
		SELECT 
			n.id, 
			n.title, 
			COALESCE(n.poster_url, ''), 
			COALESCE(ms.read_count, 0)
		FROM top_ids t
		JOIN novela n ON n.id = t.novela_id
		CROSS JOIN months m
		LEFT JOIN monthly_stats ms 
			ON ms.novela_id = t.novela_id 
			AND ms.month_start = m.month_start
		ORDER BY t.total_reads DESC, n.id, m.month_start
	`

	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats = make([]db.NovelaMonthlyRow, 0, limit)
	
	for rows.Next() {
		var s db.NovelaMonthlyRow
		if err := rows.Scan(
			&s.NovelaID,
			&s.Title,
			&s.PosterURL,
			&s.ReadCount,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

func (r *NovelaStatisticsRepository) CountNovelas(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM novela`
	var count int
	err := r.DB.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *NovelaStatisticsRepository) startOfPeriod(period constants.StatisticsPeriodSort) *time.Time {
	now := time.Now()
	var day int

	switch period {
		case constants.Day:
			day = now.Day()
		case constants.Week:
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			day = now.Day() - (weekday - 1)
		case constants.Month:
			day = 1
		default:
			return nil
	}

	t := time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location())
	return &t
}
