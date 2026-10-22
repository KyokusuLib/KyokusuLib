package db

import (
	"time"
)

type Novela struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	AlternativeTitles []string  `json:"alternative_titles"`
	Description       string    `json:"description"`
	Type              string    `json:"type"`
	AgeRating         string    `json:"age_rating"`
	ReleaseDate       time.Time `json:"release_date"`
	Status            string    `json:"status"`
	TranslationStatus string    `json:"translation_status"`
	PosterURL         string    `json:"poster_url"`
	Country           string    `json:"country"`
	Views             int       `json:"views"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Genres     []string       `json:"genres"`
	Categories []string       `json:"categories"`
	Authors    []NovelaAuthor `json:"authors"`
	Volumes    []NovelaVolume `json:"volumes"`

	Rating      float64 `json:"rating"`
	RatingCount int     `json:"rating_count"`
	UserRating  int     `json:"user_rating"`

	Bookmark      *BookmarkCategory `json:"bookmark"`
	BookmarkCount int               `json:"bookmark_count"`

	HasLiked  bool `json:"has_liked"`
	LikeCount int  `json:"like_count"`

	LastReadChapter *LastReadChapter `json:"last_read_chapter"`
	FirstChapterID  *string          `json:"first_chapter_id"`
}

type NovelaAuthor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type NovelaVolume struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	Number        int             `json:"number"`
	Status        string          `json:"status"`
	CreatedBy     int             `json:"created_by"`
	ChaptersCount int             `json:"chapters_count"`
	Chapters      []NovelaChapter `json:"chapters"`
}

type NovelaChapter struct {
	ID        string               `json:"id"`
	VolumeID  string               `json:"volume_id"`
	Title     string               `json:"title"`
	Number    float64              `json:"number"`
	Status    string               `json:"status"`
	CreatedBy int                  `json:"created_by"`
	Content   string               `json:"content,omitempty"`
	IsRead    bool                 `json:"is_read"`
	Images    []NovelaChapterImage `json:"images"`
}

type NovelaChapterImage struct {
	ID       int    `json:"id"`
	ImageURL string `json:"image_url"`
	Caption  string `json:"caption"`
	Position int    `json:"position"`
}
type NovelaLike struct {
	NovelaID int  `json:"novela_id"`
	UserID   int  `json:"user_id"`
	HasLiked bool `json:"has_liked"`
}

type NovelaRating struct {
	NovelaID int `json:"novela_id"`
	UserID   int `json:"user_id"`
	Rating   int `json:"rating"`
}

type NovelaRatingSummary struct {
	TotalCount    int
	AverageRating float64
	Distribution  map[int]int
}

type NovelaTotalRaing struct {
	Total      int                `json:"total"`
	TotalRaing float64            `json:"total_rating"`
	Rating     []NovelaRatingItem `json:"rating"`
}

type NovelaRatingItem struct {
	Value interface{} `json:"value"`
	Count int         `json:"count"`
}

type NovelaBookmarkSummary struct {
	TotalCount   int
	Distribution map[string]int
}

type UserNovelaBookmark struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	PosterURL string  `json:"poster_url"`
	Type      string  `json:"type"`
	Rating    float64 `json:"rating"`
}

type ReadChapter struct {
	UserID         int       `json:"user_id"`
	ChapterID      string    `json:"chapter_id"`
	ScrollPosition int       `json:"scroll_position"`
	CreatedAt      time.Time `json:"created_at"`
}

type LastReadChapter struct {
	ID             string  `json:"id"`
	Number         float64 `json:"number"`
	ScrollPosition int     `json:"scroll_position"`
}
