package db

import "time"

type TgPost struct {
	ID        int64         `db:"id"`
	MessageID int64         `db:"message_id"`
	Text      string        `db:"text"`
	Images    []TgPostImage `db:"-"`
	CreatedAt time.Time     `db:"created_at"`
}

type TgPostImage struct {
	ID        int64  `db:"id"`
	PostID    int64  `db:"post_id"`
	Position  int    `db:"position"`
	ImagePath string `db:"image_path"`
}
