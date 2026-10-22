package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/middleware"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/utils/validation"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	OAuthConfigGoogle  *oauth2.Config
	OAuthConfigDiscord *oauth2.Config
	AuthService        *service.AuthService
	UserService        *service.UserService
	EmailService       *service.EmailService
	SocialService      *service.SocialsService

	JWTSecret    string
	FrontendURL  string
	BackendURL   string
	CookieSecure bool
}

func NewAuthHandler(cfg *config.Config,
	authService *service.AuthService,
	userService *service.UserService,
	emailService *service.EmailService,
	socialService *service.SocialsService,
) *AuthHandler {
	return &AuthHandler{
		OAuthConfigGoogle: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		OAuthConfigDiscord: &oauth2.Config{
			ClientID:     cfg.DiscordClientID,
			ClientSecret: cfg.DiscordClientSecret,
			RedirectURL:  cfg.DiscordRedirectURL,
			Scopes:       []string{"identify", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
		AuthService:   authService,
		EmailService:  emailService,
		SocialService: socialService,
		UserService:   userService,

		JWTSecret:    cfg.JWTSecret,
		FrontendURL:  cfg.FrontendURL,
		BackendURL:   cfg.BackendURL,
		CookieSecure: cfg.CookieSecure,
	}
}

func (h *AuthHandler) newAuthCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     "KYOKUSU_API_TOKEN",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 1, 0),
		HttpOnly: true,
		Secure:   h.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (h *AuthHandler) newLogoutCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "KYOKUSU_API_TOKEN",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
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
		Token:    token.RefreshToken,
	}

	h.completeOAuth(w, r, userDTO, db.GOOGLE_PROVIDER)
}

func (h *AuthHandler) DiscordLogin(w http.ResponseWriter, r *http.Request) {
	state, err := newOAuthState(w, h.CookieSecure)
	if err != nil {
		http.Error(w, "Failed to generate OAuth state", http.StatusInternalServerError)
		return
	}

	url := h.OAuthConfigDiscord.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) DiscordCallback(w http.ResponseWriter, r *http.Request) {
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
		Token:    token.RefreshToken,
	}

	h.completeOAuth(w, r, userDTO, db.DISCORD_PROVIDER)
}

func (h *AuthHandler) completeOAuth(w http.ResponseWriter, r *http.Request, userDTO *dto.UserDTO, provider db.SocialProvider) {
	if userDTO.Email == "" {
		http.Error(w, "Email is required but provider returned empty email", http.StatusBadRequest)
		return
	}

	userDb, err := h.AuthService.LoginUser(r.Context(), userDTO)
	if err != nil {
		http.Error(w, "Login error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	h.SocialService.LinkAccount(r.Context(), userDb.ID, provider, userDTO.Id, userDTO.Token)
	h.setCookieAndRedirect(w, r, userDb)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dto dto.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.RegisterUser(r.Context(), &dto)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	verificationLink := fmt.Sprintf("%s/api/auth/verify-email?token=%s", h.BackendURL, user.VerificationToken)
	emailReq := h.EmailService.NewEmailRequest(user.Email, verificationLink)

	if err := h.EmailService.SendEmail(*emailReq); err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
		http.Error(w, "User created, but failed to send verification email.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	err := h.AuthService.VerifyUser(r.Context(), token)
	if err != nil {
		http.Error(w, "Verification failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, h.FrontendURL+"/login?verified=true", http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dto dto.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.LoginWithPassword(r.Context(), &dto)
	if err != nil {
		if err.Error() == "email not verified" {
			http.Error(w, "Email not verified. Please check your inbox.", http.StatusForbidden)
			return
		}
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	h.setCookieAndJson(w, user)
}

func (h *AuthHandler) setCookieAndRedirect(w http.ResponseWriter, r *http.Request, user *dto.GetUserDTO) {
	tokenString, err := service.GenerateJWT(user, h.JWTSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, h.newAuthCookie(tokenString))

	http.Redirect(w, r, h.FrontendURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) setCookieAndJson(w http.ResponseWriter, user *dto.GetUserDTO) {
	tokenString, err := service.GenerateJWT(user, h.JWTSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, h.newAuthCookie(tokenString))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	user, err := h.UserService.GetUserById(userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	err := h.AuthService.UpdateStatus(r.Context(), userID, false)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, h.newLogoutCookie())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "logged out"}`))
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.ForgotPasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := validation.EmailValidate(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resetToken, err := h.AuthService.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Error processing request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resetLink := fmt.Sprintf("%s/login?token=%s", h.FrontendURL, resetToken)
	emailReq := h.EmailService.NewResetPasswordEmailRequest(req.Email, resetLink)

	if err := h.EmailService.SendEmail(*emailReq); err != nil {
		http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If an account with that email exists, we have sent a password reset link.",
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.ResetPasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "Invalid password or token", http.StatusBadRequest)
		return
	}

	err := h.AuthService.ResetPassword(r.Context(), req.Token, req.Password)
	if err != nil {
		http.Error(w, "Password reset failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been successfully reset. You can now login.",
	})
}
