package response

import (
	"net/http/httptest"
	"testing"
)

type benchNovela struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	PosterURL   string   `json:"poster_url"`
	Genres      []string `json:"genres"`
	Rating      float64  `json:"rating"`
}

func BenchmarkJSON(b *testing.B) {
	payload := benchNovela{
		ID:          7,
		Title:       "Benchmark Novel",
		Description: "A moderately long description used to approximate a real catalogue payload.",
		PosterURL:   "/uploads/novelas/posters/00000000-0000-0000-0000-000000000000.jpg",
		Genres:      []string{"fantasy", "romance", "drama", "adventure"},
		Rating:      8.7,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		JSON(w, 200, payload)
	}
}

func BenchmarkError(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		Error(w, 400, "validation failed")
	}
}
