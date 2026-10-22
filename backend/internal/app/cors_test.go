package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/rs/cors"
)

func newTestCORShandler(frontendURL string) http.Handler {
	return cors.New(corsOptions(&config.Config{FrontendURL: frontendURL})).
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
}

func TestCORSAllowsConfiguredFrontend(t *testing.T) {
	handler := newTestCORShandler("https://kyokusu.example/")

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	req.Header.Set("Origin", "https://kyokusu.example")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://kyokusu.example" {
		t.Errorf("Access-Control-Allow-Origin = %q, want https://kyokusu.example", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials = %q, want true", got)
	}
	exposed := rec.Header().Get("Access-Control-Expose-Headers")
	if !strings.Contains(exposed, "X-Next-Cursor") || !strings.Contains(exposed, "X-Has-More") {
		t.Errorf("Access-Control-Expose-Headers = %q, want cursor headers", exposed)
	}
}

func TestCORSRejectsForeignOrigin(t *testing.T) {
	handler := newTestCORShandler("https://kyokusu.example")

	for _, origin := range []string{"https://evil.example", "http://localhost:3000"} {
		t.Run(origin, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
			req.Header.Set("Origin", origin)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
				t.Errorf("Access-Control-Allow-Origin = %q, want empty", got)
			}
		})
	}
}

func TestCORSPreflight(t *testing.T) {
	handler := newTestCORShandler("https://kyokusu.example")

	req := httptest.NewRequest(http.MethodOptions, "/api/novela/1", nil)
	req.Header.Set("Origin", "https://kyokusu.example")
	req.Header.Set("Access-Control-Request-Method", "PUT")
	req.Header.Set("Access-Control-Request-Headers", "content-type")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, "PUT") {
		t.Errorf("Access-Control-Allow-Methods = %q, want it to contain PUT", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(strings.ToLower(got), "content-type") {
		t.Errorf("Access-Control-Allow-Headers = %q, want it to contain content-type", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://kyokusu.example" {
		t.Errorf("Access-Control-Allow-Origin = %q, want https://kyokusu.example", got)
	}
}
