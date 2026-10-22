package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lanxre/kyokusulib/internal/models/dto"
)

func TestNewAuthCookieFlags(t *testing.T) {
	tests := []struct {
		name   string
		secure bool
	}{
		{name: "development", secure: false},
		{name: "production", secure: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{CookieSecure: tt.secure}
			cookie := h.newAuthCookie("token-value")

			if cookie.Name != "KYOKUSU_API_TOKEN" {
				t.Errorf("Name = %q, want KYOKUSU_API_TOKEN", cookie.Name)
			}
			if cookie.Value != "token-value" {
				t.Errorf("Value = %q, want token-value", cookie.Value)
			}
			if cookie.Path != "/" {
				t.Errorf("Path = %q, want /", cookie.Path)
			}
			if !cookie.HttpOnly {
				t.Error("HttpOnly must be true")
			}
			if cookie.Secure != tt.secure {
				t.Errorf("Secure = %v, want %v", cookie.Secure, tt.secure)
			}
			if cookie.SameSite != http.SameSiteLaxMode {
				t.Errorf("SameSite = %v, want Lax", cookie.SameSite)
			}
			if cookie.Expires.IsZero() {
				t.Error("Expires must be set")
			}
		})
	}
}

func TestNewLogoutCookieClearsSession(t *testing.T) {
	h := &AuthHandler{CookieSecure: true}
	cookie := h.newLogoutCookie()

	if cookie.Name != "KYOKUSU_API_TOKEN" {
		t.Errorf("Name = %q, want KYOKUSU_API_TOKEN", cookie.Name)
	}
	if cookie.Value != "" {
		t.Errorf("Value = %q, want empty", cookie.Value)
	}
	if cookie.Path != "/" {
		t.Errorf("Path = %q, want /", cookie.Path)
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

func TestSetCookieAndJsonWritesSecureCookie(t *testing.T) {
	h := &AuthHandler{
		JWTSecret:    strings.Repeat("j", 32),
		CookieSecure: true,
	}

	rec := httptest.NewRecorder()
	h.setCookieAndJson(rec, &dto.GetUserDTO{ID: 1, Email: "user@example.com", Role: "user"})

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "KYOKUSU_API_TOKEN" {
		t.Errorf("Name = %q, want KYOKUSU_API_TOKEN", cookie.Name)
	}
	if cookie.Value == "" {
		t.Error("Value is empty, want JWT")
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
	if cookie.Path != "/" {
		t.Errorf("Path = %q, want /", cookie.Path)
	}
}
