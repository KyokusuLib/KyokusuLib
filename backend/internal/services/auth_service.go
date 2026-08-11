package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
)

type AuthService struct {
	Repo            *repository.UserRepository
	UserProfileRepo *repository.UserProfileRepository
}

func NewAuthService(repo *repository.UserRepository, userProfileRepo *repository.UserProfileRepository) *AuthService {
	return &AuthService{
		Repo:            repo,
		UserProfileRepo: userProfileRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, input *dto.RegisterDTO) (*db.User, error) {
	existing, err := s.Repo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hashedPwd, err := s.hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken()
	if err != nil {
		return nil, err
	}
	
	expiresAt := time.Now().Add(24 * time.Hour)

	user := &db.User{
		Email:                      input.Email,
		Name:                       input.Username,
		PasswordHash:               hashedPwd,
		Role:                       "user",
		Status:                     "online",
		IsVerified:                 false,
		VerificationToken:          token,
		VerificationTokenExpiresAt: &expiresAt,
		Gender:                     db.HIDDEN_GENERE,
	}

	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginWithPassword(ctx context.Context, input *dto.LoginDTO) (*dto.GetUserDTO, error) {
	user, err := s.Repo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return nil, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := s.Repo.UpdateStatus(user.ID, "online"); err != nil {
		log.Printf("failed to update status: %v", err)
	}
	user.Status = "online"

	return s.enrichUserDTO(ctx, user)
}

func (s *AuthService) LoginUser(ctx context.Context, gUser *dto.UserDTO) (*dto.GetUserDTO, error) {
	user, err := s.Repo.GetByEmail(gUser.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &db.User{
			Email:      gUser.Email,
			Name:       gUser.Username,
			Picture:    gUser.Avatar,
			Role:       "user",
			Status:     "online",
			IsVerified: true,
		}
		if err := s.Repo.Create(user); err != nil {
			return nil, err
		}

		return s.toUserDTO(user, nil, nil, 0, 0), nil
	}

	if !user.IsVerified {
		if err := s.Repo.MarkUserVerified(user.ID); err != nil {
			log.Printf("failed to verify user: %v", err)
		}
		user.IsVerified = true
	}

	if user.Picture == "" && gUser.Avatar != "" {
		if err := s.Repo.UpdateAvatar(user.ID, gUser.Avatar); err != nil {
			log.Printf("failed to update avatar: %v", err)
		}
		user.Picture = gUser.Avatar
	}

	if err := s.Repo.UpdateStatus(user.ID, "online"); err != nil {
		log.Printf("failed to update status: %v", err)
	}
	user.Status = "online"

	return s.enrichUserDTO(ctx, user)
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	token, err := s.generateToken()
	if err != nil {
		return "", err
	}
	
	expiresAt := time.Now().Add(1 * time.Hour)

	user.ResetToken = token
	user.ResetTokenExpiresAt = &expiresAt

	if err := s.Repo.Update(user); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.Repo.GetByResetToken(token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid or expired token")
	}

	if user.ResetTokenExpiresAt != nil && user.ResetTokenExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	hashed, err := s.hashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashed
	user.ResetToken = ""
	user.ResetTokenExpiresAt = nil

	return s.Repo.Update(user)
}

func (s *AuthService) VerifyUser(ctx context.Context, token string) error {
	user, err := s.Repo.GetByVerificationToken(token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid token")
	}

	if user.VerificationTokenExpiresAt != nil && user.VerificationTokenExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	return s.Repo.MarkUserVerified(user.ID)
}

func (s *AuthService) UpdateStatus(ctx context.Context, userID int, isOnline bool) error {
	status := "offline"
	if isOnline {
		status = "online"
	}
	
	return s.Repo.UpdateStatus(userID, status)
}

func (s *AuthService) enrichUserDTO(ctx context.Context, user *db.User) (*dto.GetUserDTO, error) {
	g, gctx := errgroup.WithContext(ctx)

	var userTags []dto.UserTagDTO
	g.Go(func() error {
		tagsDB, err := s.Repo.GetUserTags(gctx, user.ID)
		if err != nil {
			return err
		}
		userTags = toUserTagDTOs(tagsDB)
		return nil
	})

	var userLevel *db.UserLevel
	g.Go(func() error {
		level, err := s.UserProfileRepo.GetUserLevel(gctx, user.ID)
		if err != nil {
			log.Printf("failed to get user level: %v", err)
		}
		userLevel = level
		return nil
	})

	var totalComments, readChapters int
	g.Go(func() error {
		comments, chapters, err := s.Repo.GetUserStats(gctx, user.ID)
		if err != nil {
			log.Printf("failed to get user stats: %v", err)
		}
		totalComments, readChapters = comments, chapters
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return s.toUserDTO(user, userTags, userLevel, totalComments, readChapters), nil
}

func (s AuthService) toUserDTO(user *db.User, tags []dto.UserTagDTO, level *db.UserLevel, totalComments int, readChapters int) *dto.GetUserDTO {
	res := &dto.GetUserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Picture:   user.Picture,
		Role:      user.Role,
		Status:    user.Status,
		About:     user.About,
		Birthday:  user.Birthday,
		Gender:    string(user.Gender),
		IsPublic:  user.IsPublic,
		LastLogin: user.LastLogin,
		CreateAt:  user.CreateAt,
		Banner:    user.Banner,
		ActiveTag: user.Tag,
		AllTags:   tags,
		Settings: dto.PublicUserSettingsDTO{
			IsShowTag: user.IsShowTag,
			IsShowBookmark: user.IsShowBookmark,
		},
		UserStats: dto.UserStatsDTO{
			TotalComments: totalComments,
			ReadChapters:  readChapters,
		},
	}

	if level != nil {
		res.UserLevel = dto.UserLevelDTO{
			Level:      level.Level,
			Experience: level.Experience,
			LevelTitle: level.LevelTitle,
			XPForNext:  level.XPForNext,
		}
	}

	return res
}

func (s *AuthService) hashPassword(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(b), err
}

func (s *AuthService) generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *AuthService) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("User cleanup worker stopped")
				return
			case <-ticker.C:
				log.Println("Running cleanup task...")
				if err := s.Repo.DeleteExpiredUnverified(); err != nil {
					log.Printf("Error cleaning up users: %v", err)
				}
			}
		}
	}()
}