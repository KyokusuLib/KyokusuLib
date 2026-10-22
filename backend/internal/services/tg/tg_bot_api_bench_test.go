package tg

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/lanxre/kyokusulib/internal/config"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newBenchBotAPI(body string, status int) *TgBotAPI {
	api := NewTgBotAPI(
		&config.Config{TgBotToken: "bench-token", TgChannelUsername: "@bench"},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	api.client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    r,
		}, nil
	})
	return api
}

func BenchmarkDeleteFromTelegramOK(b *testing.B) {
	api := newBenchBotAPI(`{"ok":true}`, http.StatusOK)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := api.DeleteFromTelegram(context.Background(), 42); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeleteFromTelegramNotFound(b *testing.B) {
	api := newBenchBotAPI(`{"ok":false,"error_code":400,"description":"Bad Request: message to delete not found"}`, http.StatusBadRequest)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := api.DeleteFromTelegram(context.Background(), 42); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeleteFromTelegramPermanent(b *testing.B) {
	api := newBenchBotAPI(`{"ok":false,"error_code":400,"description":"Bad Request: message can't be deleted for everyone"}`, http.StatusBadRequest)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = api.DeleteFromTelegram(context.Background(), 42)
	}
}
