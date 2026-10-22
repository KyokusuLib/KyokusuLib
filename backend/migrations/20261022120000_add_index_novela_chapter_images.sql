-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_novela_chapter_images_chapter_id_position
    ON novela_chapter_images (chapter_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_novela_chapter_images_chapter_id_position;
-- +goose StatementEnd
