-- +goose Up
-- +goose StatementBegin
CREATE TABLE tg_post_images (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES tg_posts(id) ON DELETE CASCADE,
    position INT NOT NULL DEFAULT 0,
    image_path TEXT NOT NULL
);

CREATE INDEX idx_tg_post_images_post_id ON tg_post_images(post_id);

INSERT INTO tg_post_images (post_id, position, image_path)
SELECT id, 0, image_path FROM tg_posts WHERE image_path <> '';

ALTER TABLE tg_posts DROP COLUMN image_path;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tg_posts ADD COLUMN image_path TEXT NOT NULL DEFAULT '';

UPDATE tg_posts
SET image_path = COALESCE(
    (SELECT image_path FROM tg_post_images WHERE post_id = tg_posts.id ORDER BY position ASC LIMIT 1),
    ''
);

DROP TABLE IF EXISTS tg_post_images;
-- +goose StatementEnd
