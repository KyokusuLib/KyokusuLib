package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOAuthStateSetsCookie(t *testing.T) {
	rec := httptest.NewRecorder()
	state, err := newOAuthState(rec, true)
	if err != nil {
		t.Fatalf("newOAuthState() error = %v", err)
	}
	if state == "" {
		t.Fatal("state is empty")
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != oauthStateCookie {
		t.Errorf("Name = %q, want %q", cookie.Name, oauthStateCookie)
	}
	if cookie.Value != state {
		t.Errorf("Value = %q, want %q", cookie.Value, state)
	}
	if cookie.Path != "/" {
		t.Errorf("Path = %q, want /", cookie.Path)
	}
	if cookie.MaxAge != int(oauthStateTTL.Seconds()) {
		t.Errorf("MaxAge = %d, want %d", cookie.MaxAge, int(oauthStateTTL.Seconds()))
	}
	if !cookie.HttpOnly {
		t.Error("HttpOnly must be true")
	}
	if !cookie.Secure {
		t.Error("Secure must be true")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("SameSite = %v, want Lax", cookie.SameSite)
	}
}

func TestNewOAuthStateIsRandom(t *testing.T) {
	first, err := newOAuthState(httptest.NewRecorder(), false)
	if err != nil {
		t.Fatalf("newOAuthState() error = %v", err)
	}
	second, err := newOAuthState(httptest.NewRecorder(), false)
	if err != nil {
		t.Fatalf("newOAuthState() error = %v", err)
	}
	if first == second {
		t.Error("two OAuth states must differ")
	}
}

func TestValidOAuthState(t *testing.T) {
	tests := []struct {
		name      string
		state     string
		cookie    string
		hasCookie bool
		want      bool
	}{
		{name: "match", state: "expected-state", cookie: "expected-state", hasCookie: true, want: true},
		{name: "mismatch", state: "expected-state", cookie: "other-state", hasCookie: true, want: false},
		{name: "missing cookie", state: "expected-state", hasCookie: false, want: false},
		{name: "empty state", state: "", cookie: "expected-state", hasCookie: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/callback", nil)
			if tt.hasCookie {
				req.AddCookie(&http.Cookie{Name: oauthStateCookie, Value: tt.cookie})
			}

			if got := validOAuthState(req, tt.state); got != tt.want {
				t.Errorf("validOAuthState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearOAuthState(t *testing.T) {
	rec := httptest.NewRecorder()
	clearOAuthState(rec, true)

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Value != "" {
		t.Errorf("Value = %q, want empty", cookie.Value)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("MaxAge = %d, want -1", cookie.MaxAge)
	}
	if !cookie.HttpOnly {
		t.Error("HttpOnly must be true")
	}
	if !cookie.Secure {
		t.Error("Secure must be true")
	}
}

func TestGoogleCallbackRejectsInvalidState(t *testing.T) {
	h := &AuthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/auth/google/callback?state=wrong", nil)
	req.AddCookie(&http.Cookie{Name: oauthStateCookie, Value: "expected"})
	rec := httptest.NewRecorder()

	h.GoogleCallback(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
