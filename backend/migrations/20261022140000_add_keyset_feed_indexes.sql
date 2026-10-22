-- +goose Up
-- +goose StatementBegin
-- Keyset pagination orders feeds by (time, id) descending. The id tie-breaker
-- keeps pages stable when many rows share the same timestamp.
CREATE INDEX idx_tg_posts_created_at_id
    ON tg_posts (created_at DESC, id DESC);

CREATE INDEX idx_notifications_user_id_created_at_id
    ON notifications (user_id, created_at DESC, id DESC);

CREATE INDEX idx_user_activities_user_id_timestamp_id
    ON user_activities (user_id, "timestamp" DESC, id DESC);

-- The wider indexes above cover the old prefix-only indexes.
DROP INDEX IF EXISTS idx_tg_posts_created_at;
DROP INDEX IF EXISTS idx_notifications_user_id_created_at;
DROP INDEX IF EXISTS idx_user_activities_user_id_timestamp;
-- No global activity feed exists; only per-user keyset queries use this table.
DROP INDEX IF EXISTS idx_user_activities_timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_activities_timestamp;

CREATE INDEX idx_tg_posts_created_at
    ON tg_posts (created_at DESC);

CREATE INDEX idx_notifications_user_id_created_at
    ON notifications (user_id, created_at DESC);

CREATE INDEX idx_user_activities_user_id_timestamp
    ON user_activities (user_id, "timestamp" DESC);

CREATE INDEX idx_user_activities_timestamp
    ON user_activities ("timestamp" DESC);

DROP INDEX IF EXISTS idx_user_activities_user_id_timestamp_id;
DROP INDEX IF EXISTS idx_notifications_user_id_created_at_id;
DROP INDEX IF EXISTS idx_tg_posts_created_at_id;
-- +goose StatementEnd
