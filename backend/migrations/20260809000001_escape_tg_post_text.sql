-- +goose Up
-- +goose StatementBegin
-- Existing rows hold plain Telegram text; since the bot now stores HTML,
-- escape them so they render safely as HTML on the frontend.
UPDATE tg_posts
SET text = replace(replace(replace(text, '&', '&amp;'), '<', '&lt;'), '>', '&gt;');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Irreversible: unescaping plain text from HTML is not reliable.
SELECT 1;
-- +goose StatementEnd
