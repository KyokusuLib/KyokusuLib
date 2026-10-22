-- +goose Up
-- +goose StatementBegin
-- On managed Postgres the extension may need to be enabled by an administrator
-- before this migration can run under the application role.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Concatenates the title and alternative titles so one trigram index can
-- serve both fields. array_to_string() is marked STABLE, so wrap it in an
-- IMMUTABLE function to make it usable in an index expression.
CREATE OR REPLACE FUNCTION novela_search_text(title text, alternative_titles text[])
RETURNS text
LANGUAGE sql
IMMUTABLE
PARALLEL SAFE
AS $$
    SELECT title || ' ' || COALESCE(array_to_string(alternative_titles, ' '), '')
$$;

-- Notifications are always read per user, newest first.
CREATE INDEX idx_notifications_user_id_created_at
    ON notifications (user_id, created_at DESC);

-- Activity feed: per user, newest first.
CREATE INDEX idx_user_activities_user_id_timestamp
    ON user_activities (user_id, "timestamp" DESC);

-- Bookmarks are listed per user and category, newest first.
CREATE INDEX idx_user_novela_bookmarks_user_category_updated
    ON user_novela_bookmarks (user_id, category_id, updated_at DESC);

-- Reverse lookups for catalog filters by genre/category.
CREATE INDEX idx_novela_genres_genre_id_novela_id
    ON novela_genres (genre_id, novela_id);
CREATE INDEX idx_novela_categories_category_id_novela_id
    ON novela_categories (category_id, novela_id);

-- Default catalog ordering and the "trending" sort.
CREATE INDEX idx_novela_created_at
    ON novela (created_at DESC);
CREATE INDEX idx_novela_views
    ON novela (views DESC);

-- last_chapter_at and chapter-count subqueries only look at approved chapters.
CREATE INDEX idx_novela_chapters_volume_created_approved
    ON novela_chapters (novela_volume_id, created_at DESC)
    WHERE status = 'approved';

-- Substring search over title + alternative titles.
CREATE INDEX idx_novela_search_trgm
    ON novela USING gin (novela_search_text(title::text, alternative_titles) gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_novela_search_trgm;
DROP INDEX IF EXISTS idx_novela_chapters_volume_created_approved;
DROP INDEX IF EXISTS idx_novela_views;
DROP INDEX IF EXISTS idx_novela_created_at;
DROP INDEX IF EXISTS idx_novela_categories_category_id_novela_id;
DROP INDEX IF EXISTS idx_novela_genres_genre_id_novela_id;
DROP INDEX IF EXISTS idx_user_novela_bookmarks_user_category_updated;
DROP INDEX IF EXISTS idx_user_activities_user_id_timestamp;
DROP INDEX IF EXISTS idx_notifications_user_id_created_at;
DROP FUNCTION IF EXISTS novela_search_text(text, text[]);
-- +goose StatementEnd
