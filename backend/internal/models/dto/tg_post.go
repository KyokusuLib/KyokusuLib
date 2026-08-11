package dto

import "time"

type TgPost struct {
	ID        int64     `json:"id"`
	MessageID int64     `json:"messageId"`
	Text      string    `json:"text"`
	ImageURLs []string  `json:"imageUrls"`
	CreatedAt time.Time `json:"createdAt"`
}

type DeleteMessageResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}
