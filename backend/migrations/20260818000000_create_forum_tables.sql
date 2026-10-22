-- +goose Up
-- +goose StatementBegin

CREATE TABLE forum_posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    view_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forum_posts_user_id ON forum_posts(user_id);
CREATE INDEX idx_forum_posts_created_at ON forum_posts(created_at DESC);

CREATE TABLE forum_post_images (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES forum_posts(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forum_post_images_post_id ON forum_post_images(post_id);

CREATE TABLE forum_comments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES forum_posts(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id INT REFERENCES forum_comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forum_comments_post_id ON forum_comments(post_id);
CREATE INDEX idx_forum_comments_parent_id ON forum_comments(parent_id);

CREATE TABLE forum_post_likes (
    post_id INT NOT NULL REFERENCES forum_posts(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, user_id)
);

CREATE TABLE forum_comment_likes (
    comment_id INT NOT NULL REFERENCES forum_comments(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id, user_id)
);

CREATE TABLE forum_report_reasons (
    id SERIAL PRIMARY KEY,
    reason VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE forum_post_reports (
    post_id INT NOT NULL REFERENCES forum_posts(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason_id INT NOT NULL REFERENCES forum_report_reasons(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, user_id)
);

CREATE TABLE forum_comment_reports (
    comment_id INT NOT NULL REFERENCES forum_comments(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason_id INT NOT NULL REFERENCES forum_report_reasons(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id, user_id)
);

INSERT INTO forum_report_reasons (reason) VALUES
    ('Нарушение правил'),
    ('Спам'),
    ('Другое');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS forum_comment_reports;
DROP TABLE IF EXISTS forum_post_reports;
DROP TABLE IF EXISTS forum_report_reasons;
DROP TABLE IF EXISTS forum_comment_likes;
DROP TABLE IF EXISTS forum_post_likes;
DROP TABLE IF EXISTS forum_comments;
DROP TABLE IF EXISTS forum_post_images;
DROP TABLE IF EXISTS forum_posts;
-- +goose StatementEnd
