package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/middleware"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	service "github.com/lanxre/kyokusulib/internal/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type SocialNetworkHandler struct {
	Service            *service.SocialsService
	OAuthConfigGoogle  *oauth2.Config
	OAuthConfigDiscord *oauth2.Config
	FrontendURL        string
	CookieSecure       bool
}

func NewSocialNetworkHandler(cfg *config.Config, service *service.SocialsService) *SocialNetworkHandler {
	return &SocialNetworkHandler{
		Service: service,
		OAuthConfigGoogle: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleSocialRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		OAuthConfigDiscord: &oauth2.Config{
			ClientID:     cfg.DiscordClientID,
			ClientSecret: cfg.DiscordClientSecret,
			RedirectURL:  cfg.DiscordSocialRedirectURL,
			Scopes:       []string{"identify", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
		FrontendURL:  cfg.FrontendURL,
		CookieSecure: cfg.CookieSecure,
	}
}

func (h *SocialNetworkHandler) GetSocials(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	socials, err := h.Service.GetUserSocials(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(socials)
}

func (h *SocialNetworkHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := newOAuthState(w, h.CookieSecure)
	if err != nil {
		http.Error(w, "Failed to generate OAuth state", http.StatusInternalServerError)
		return
	}

	url := h.OAuthConfigGoogle.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *SocialNetworkHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	state := r.FormValue("state")
	if !validOAuthState(r, state) {
		clearOAuthState(w, h.CookieSecure)
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	clearOAuthState(w, h.CookieSecure)

	code := r.FormValue("code")
	token, err := h.OAuthConfigGoogle.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := h.OAuthConfigGoogle.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var user dto.GoogleUser
	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	userDTO := &dto.UserDTO{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.Name,
		Avatar:   user.Picture,
	}

	h.completeOAuth(w, r, userID, userDTO, db.GOOGLE_PROVIDER, token.RefreshToken)
}

func (h *SocialNetworkHandler) DiscordLogin(w http.ResponseWriter, r *http.Request) {
	state, err := newOAuthState(w, h.CookieSecure)
	if err != nil {
		http.Error(w, "Failed to generate OAuth state", http.StatusInternalServerError)
		return
	}

	url := h.OAuthConfigDiscord.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *SocialNetworkHandler) DiscordCallback(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}
	state := r.FormValue("state")
	if !validOAuthState(r, state) {
		clearOAuthState(w, h.CookieSecure)
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	clearOAuthState(w, h.CookieSecure)

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := h.OAuthConfigDiscord.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := h.OAuthConfigDiscord.Client(context.Background(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		http.Error(w, "Failed to get discord user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var dUser dto.DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&dUser); err != nil {
		http.Error(w, "Failed to parse discord response", http.StatusInternalServerError)
		return
	}

	avatarURL := ""
	if dUser.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", dUser.ID, dUser.Avatar)
	} else {
		avatarURL = "https://cdn.discordapp.com/embed/avatars/0.png"
	}

	displayName := dUser.GlobalName
	if displayName == "" {
		displayName = dUser.Username
	}

	userDTO := &dto.UserDTO{
		Id:       dUser.ID,
		Email:    dUser.Email,
		Username: displayName,
		Avatar:   avatarURL,
	}

	h.completeOAuth(w, r, userID, userDTO, db.DISCORD_PROVIDER, token.RefreshToken)
}

func (h *SocialNetworkHandler) UnlinkGoogle(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	socials, err := h.Service.GetUserSocials(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user socials", http.StatusInternalServerError)
		return
	}

	if socials.GoogleRefreshToken != "" {
		revokeURL := "https://oauth2.googleapis.com/revoke"
		data := url.Values{}
		data.Set("token", socials.GoogleRefreshToken)
		http.PostForm(revokeURL, data)
	}

	if err := h.Service.UnlinkAccount(r.Context(), userID, db.GOOGLE_PROVIDER); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Google account unlinked"})
}

func (h *SocialNetworkHandler) UnlinkDiscord(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	socials, err := h.Service.GetUserSocials(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user socials", http.StatusInternalServerError)
		return
	}

	if socials.DiscordRefreshToken != "" {
		data := url.Values{}
		data.Set("token", socials.DiscordRefreshToken)
		data.Set("token_type_hint", "refresh_token")

		req, _ := http.NewRequest("POST", "https://discord.com/api/oauth2/token/revoke", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(h.OAuthConfigDiscord.ClientID, h.OAuthConfigDiscord.ClientSecret)

		client := &http.Client{}
		client.Do(req)
	}

	if err := h.Service.UnlinkAccount(r.Context(), userID, db.DISCORD_PROVIDER); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Discord account unlinked"})
}

func (h *SocialNetworkHandler) completeOAuth(w http.ResponseWriter, r *http.Request, userID int, userDTO *dto.UserDTO, provider db.SocialProvider, refreshToken string) {
	if userDTO.Email == "" {
		http.Error(w, "Email is required but provider returned empty email", http.StatusBadRequest)
		return
	}

	h.Service.LinkAccount(r.Context(), userID, provider, userDTO.Id, refreshToken)
	http.Redirect(w, r, h.FrontendURL+"/profile/settings", http.StatusTemporaryRedirect)
}
