package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lanxre/kyokusulib/internal/models/dto"
	service "github.com/lanxre/kyokusulib/internal/services"
)

const benchJWTSecret = "bench-secret"

func benchToken(b *testing.B) string {
	b.Helper()
	token, err := service.GenerateJWT(&dto.GetUserDTO{
		ID:    1,
		Email: "bench@example.com",
		Role:  "user",
	}, benchJWTSecret)
	if err != nil {
		b.Fatal(err)
	}
	return token
}

func BenchmarkAuthMiddleware(b *testing.B) {
	token := benchToken(b)
	handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {}, benchJWTSecret)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/novela", nil)
		req.AddCookie(&http.Cookie{Name: "KYOKUSU_API_TOKEN", Value: token})
		rec := httptest.NewRecorder()
		handler(rec, req)
	}
}

func BenchmarkDefaultMiddleware(b *testing.B) {
	token := benchToken(b)
	handler := DefaultMiddleware(func(w http.ResponseWriter, r *http.Request) {}, benchJWTSecret)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/novela", nil)
		req.AddCookie(&http.Cookie{Name: "KYOKUSU_API_TOKEN", Value: token})
		rec := httptest.NewRecorder()
		handler(rec, req)
	}
}

func BenchmarkDefaultMiddlewareAnonymous(b *testing.B) {
	handler := DefaultMiddleware(func(w http.ResponseWriter, r *http.Request) {}, benchJWTSecret)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/novela", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)
	}
}
