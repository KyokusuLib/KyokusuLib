package service

import (
	"context"
	"errors"
	"time"

	"github.com/lanxre/kyokusulib/internal/constants"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"golang.org/x/sync/errgroup"
)

type UserService struct {
	Repo *repository.UserRepository
	UserProfileRepo *repository.UserProfileRepository
}

func NewUserService(repo *repository.UserRepository, userProfileRepo *repository.UserProfileRepository) *UserService {
	return &UserService{Repo: repo, UserProfileRepo: userProfileRepo}
}

func (s *UserService) GetUsers(ctx context.Context, search string, limit int, offset int) ([]*dto.GetUserDTO, error) {
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	
	usersDb, err := s.Repo.GetUsers(ctx, search, limit, offset)
	if err != nil {
		return nil, err
	}

	userIDs := make([]int, len(usersDb))
	for i, user := range usersDb {
		userIDs[i] = user.ID
	}

	levels, err := s.UserProfileRepo.GetUserLevelsBatch(context.Background(), userIDs)
	if err != nil {
		return nil, err
	}
	tags, err := s.Repo.GetUserTagsBatch(context.Background(), userIDs)
	if err != nil {
		return nil, err
	}
	stats, err := s.Repo.GetUserStatsBatch(context.Background(), userIDs)
	if err != nil {
		return nil, err
	}

	usersDto := make([]*dto.GetUserDTO, len(usersDb))
	for i, user := range usersDb {
		userLevel := levels[user.ID]
		if userLevel == nil {
			userLevel = &db.UserLevel{}
		}
		userStats := stats[user.ID]
		usersDto[i] = toUserDTO(user, userLevel, toUserTagDTOs(tags[user.ID]), dto.PublicUserSettingsDTO{
			IsShowTag:      user.IsShowTag,
			IsShowBookmark: user.IsShowBookmark,
		}, dto.UserStatsDTO{
			TotalComments: userStats.TotalComments,
			ReadChapters:  userStats.ReadChapters,
		})
	}
	return usersDto, nil
}

func (s *UserService) GetUserById(userId int) (*dto.GetUserDTO, error) {
	userDb, err := s.Repo.GetByID(userId)
	if err != nil || userDb == nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(context.Background())

	var userTags []dto.UserTagDTO
	g.Go(func() error {
		tags, err := s.GetUserTags(userId)
		if err != nil {
			return err
		}
		userTags = tags
		return nil
	})

	var userLevel *db.UserLevel
	g.Go(func() error {
		level, err := s.UserProfileRepo.GetUserLevel(ctx, userId)
		if err != nil {
			return err
		}
		userLevel = level
		return nil
	})

	var totalComments, readChapters int
	g.Go(func() error {
		comments, chapters, err := s.Repo.GetUserStats(ctx, userId)
		if err != nil {
			return err
		}
		totalComments, readChapters = comments, chapters
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return toUserDTO(userDb, userLevel, userTags, dto.PublicUserSettingsDTO{
		IsShowTag:      userDb.IsShowTag,
		IsShowBookmark: userDb.IsShowBookmark,
	}, dto.UserStatsDTO{
		TotalComments: totalComments,
		ReadChapters:  readChapters,
	}), nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, userID int, dto dto.UpdateUserStatusDTO) error {
    lastActive := time.Now()
    if dto.LastActive != 0 {
        lastActive = time.UnixMilli(dto.LastActive)
    }

    status := constants.UserStatus(dto.Status)
    if !status.IsValid() {
        return errors.New("Invalid user status")
    }

    return s.Repo.UpdateDtoStatus(ctx, userID, status, lastActive)
}

func (s *UserService) UpdateUserTag(ctx context.Context, userId int, dto dto.UpdateUserTagDTO) error {
	return s.Repo.UpdateUserTag(ctx, userId, dto.ID)
}

func (s *UserService) UpdateUser(ctx context.Context, userID int, input dto.UpdateUserDTO) error {
	exists, err := s.Repo.IsExist(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	var birthday *time.Time
	if input.Birthday != "" {
		parsed, err := time.Parse("2006-01-02", input.Birthday)
		if err != nil {
			return errors.New("invalid birthday format")
		}
		birthday = &parsed
	}

	return s.Repo.UpdateUser(ctx, userID, repository.UpdateUserParams{
		Name:           input.Name,
		About:          input.About,
		Gender:         input.Gender,
		Birthday:       birthday,
		IsPublic:       input.IsPublic,
		Role:           input.Role,
		Status:         input.Status,
		IsShowTag:      input.IsShowTag,
		IsShowBookmark: input.IsShowBookmark,
	})
}

func (s *UserService) GetUserTags(userId int) ([]dto.UserTagDTO, error) {
	tags, err := s.Repo.GetUserTags(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	return toUserTagDTOs(tags), nil
}

func toUserTagDTOs(tags []*db.UserTag) []dto.UserTagDTO {
	userTags := make([]dto.UserTagDTO, len(tags))
	for i, tag := range tags {
		userTags[i] = dto.UserTagDTO{
			ID:  tag.TagID,
			Tag: tag.Tag,
		}
	}
	return userTags
}

func (s *UserService) DeleteUserById(ctx context.Context, userID int) error {
    exists, err := s.Repo.IsExist(ctx, userID)
    if err != nil {
        return err
    }

    if !exists {
        return errors.New("user not found")
    }

    return s.Repo.Delete(ctx, userID)
}

func toUserDTO(user *db.User, level *db.UserLevel, tags []dto.UserTagDTO, settings dto.PublicUserSettingsDTO, stats dto.UserStatsDTO) *dto.GetUserDTO {
	return &dto.GetUserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Picture:   user.Picture,
		Banner:    user.Banner,
		Role:      user.Role,
		Status:    user.Status,
		About:     user.About,
		Birthday:  user.Birthday,
		Gender:    string(user.Gender),
		IsPublic:  user.IsPublic,
		LastLogin: user.LastLogin,
		CreateAt:  user.CreateAt,
		ActiveTag: user.Tag,
		AllTags:   tags,
		Settings:  settings,
		UserStats: stats,
		UserLevel: dto.UserLevelDTO{
			Level:       level.Level,
			Experience:  level.Experience,
			LevelTitle:  level.LevelTitle,
			XPForNext:   level.XPForNext,
		},
	}
}