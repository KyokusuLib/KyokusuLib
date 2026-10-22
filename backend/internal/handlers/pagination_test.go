package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/utils/cursor"
)

func TestParseFeedPagination(t *testing.T) {
	validCursor := cursor.Encode(time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC), 42)

	tests := []struct {
		name        string
		query       string
		defaultSize int
		wantLimit   int
		wantCursor  bool
		wantErr     bool
	}{
		{name: "defaults", defaultSize: 20, wantLimit: 20},
		{name: "explicit limit", query: "?limit=5", defaultSize: 20, wantLimit: 5},
		{name: "zero limit uses default", query: "?limit=0", defaultSize: 20, wantLimit: 20},
		{name: "negative limit", query: "?limit=-1", defaultSize: 20, wantErr: true},
		{name: "limit too large", query: "?limit=101", defaultSize: 20, wantErr: true},
		{name: "zero offset is accepted", query: "?offset=0", defaultSize: 20, wantLimit: 20},
		{name: "deep offset is rejected", query: "?offset=10", defaultSize: 20, wantErr: true},
		{name: "valid cursor", query: "?cursor=" + validCursor, defaultSize: 20, wantLimit: 20, wantCursor: true},
		{name: "invalid cursor", query: "?cursor=not-a-cursor", defaultSize: 20, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/feed"+tt.query, nil)
			limit, after, err := parseFeedPagination(req, tt.defaultSize)

			if tt.wantErr {
				if err == nil {
					t.Fatal("parseFeedPagination() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseFeedPagination() error = %v", err)
			}
			if limit != tt.wantLimit {
				t.Errorf("limit = %d, want %d", limit, tt.wantLimit)
			}
			if (after != nil) != tt.wantCursor {
				t.Errorf("after = %v, want cursor %v", after, tt.wantCursor)
			}
			if after != nil && after.ID != 42 {
				t.Errorf("cursor ID = %d, want 42", after.ID)
			}
		})
	}
}

func TestSetCursorHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	setCursorHeaders(rec, dto.CursorPage[int]{
		Items:      []int{1, 2},
		NextCursor: "next",
		HasMore:    true,
	})

	if got := rec.Header().Get("X-Next-Cursor"); got != "next" {
		t.Errorf("X-Next-Cursor = %q, want next", got)
	}
	if got := rec.Header().Get("X-Has-More"); got != "true" {
		t.Errorf("X-Has-More = %q, want true", got)
	}
}
