package db

import "time"

type ShortNovelaInfo struct {
	NovelaID          int
	Title             string
	Description       string
	PosterURL         string
	ReleaseDate       time.Time
}

type TotalNovelaStatistics struct {
	Novela             ShortNovelaInfo
	
	TotalBookmarkCount int
	TotalReadCount     int
	TotalRatingCount   int
	TotalCommentCount  int
	TotalVolumeCount   int
	TotalChapterCount  int
	
	Rating             float64
}

type GeneralStatistics struct {
	GeneralBookmarkCount int
	GeneralReadCount     int
	GeneralRatingCount   int
	GeneralCommentCount  int
	GeneralVolumeCount   int
	GeneralChapterCount  int
}

type NovelaMonthlyRow struct {
	NovelaID  int
	Title     string
	PosterURL string

	ReadCount int
}

type UserStats struct {
	TotalComments int
	ReadChapters  int
}
