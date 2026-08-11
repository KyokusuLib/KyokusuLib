package service

import (
	"context"
	"errors"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
)

type UserExperienceService struct {
	UserExpRepo *repository.UserExperienceRepository
}

func NewUserExperianceService(userExpRepo *repository.UserExperienceRepository) *UserExperienceService {
	return &UserExperienceService{
		UserExpRepo: userExpRepo,
	}
}

func (s *UserExperienceService) GetLevelDefinitions (ctx context.Context) ([]dto.LevelDefinitions, error) {
	defenitions, err := s.UserExpRepo.GetLevelDefinitions(ctx)
	if err != nil {
		return nil, err
	}

	levelDefinitions := make([]dto.LevelDefinitions, len(defenitions))
	for i, level := range defenitions {
		levelDefinitions[i] = s.mapToDtoLevelDefinitions(level)
	}

	return levelDefinitions, nil
}

func (s *UserExperienceService) UpdateUserLevel(ctx context.Context, dto dto.UpdateUserLevel) error {
	if dto.Level < 1 {
		return errors.New("invalid level")
	}
	if dto.Experience < 0 {
		return errors.New("invalid experience")
	}

	return s.UserExpRepo.UpdateUserLevel(
		ctx,
		dto.UserID,
		dto.Level,
		dto.Experience,
	)
}

func (s *UserExperienceService) mapToDtoLevelDefinitions(level db.UserExperienceDefinitions) dto.LevelDefinitions {
	return dto.LevelDefinitions{
		Level: level.Level,
		Title: level.Title,
		Total_XP_Required: level.Total_XP_Required,
	}
}

