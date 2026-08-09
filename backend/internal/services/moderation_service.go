package service

import (
	"context"
	"fmt"

	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
)

type ModerationService struct {
	NovelaRepo    *repository.NovelaRepository
	NovelaService *NovelaService
}

func NewModerationService(novelaRepo *repository.NovelaRepository, novelaService *NovelaService) *ModerationService {
	return &ModerationService{
		NovelaRepo:    novelaRepo,
		NovelaService: novelaService,
	}
}

func (s *ModerationService) GetPendingContent(ctx context.Context) (map[string]any, error) {
	volumes, err := s.NovelaRepo.GetPendingVolumes(ctx)
	if err != nil {
		return nil, err
	}

	chapters, err := s.NovelaRepo.GetPendingChapters(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"volumes":  volumes,
		"chapters": chapters,
	}, nil
}

func (s *ModerationService) ApproveVolume(ctx context.Context, id string) error {
	v, err := s.NovelaRepo.GetVolumeByID(ctx, id)
	if err != nil {
		return err
	}
	if v == nil {
		return fmt.Errorf("volume not found")
	}

	err = s.NovelaRepo.UpdateVolumeStatus(ctx, id, "approved")
	if err != nil {
		return err
	}

	novelaID, err := s.NovelaRepo.GetNovelaIDByVolumeID(ctx, id)
	if err == nil {
		go s.NovelaService.sendNovelaNotification(novelaID, "Обновление", "В новеллу '%s' добавлен новый том!")
	}

	return nil
}

func (s *ModerationService) RejectVolume(ctx context.Context, id string) error {
	return s.NovelaRepo.UpdateVolumeStatus(ctx, id, "rejected")
}

func (s *ModerationService) ApproveChapter(ctx context.Context, id string) error {
	ch, err := s.NovelaRepo.GetChapterByID(ctx, id)
	if err != nil {
		return err
	}
	if ch == nil {
		return fmt.Errorf("chapter not found")
	}

	err = s.NovelaRepo.UpdateChapterStatus(ctx, id, "approved")
	if err != nil {
		return err
	}

	novelaID, err := s.NovelaRepo.GetNovelaIDByChapterID(ctx, id)
	if err == nil {
		go s.NovelaService.sendNovelaNotification(novelaID, "Обновление", "В новеллу '%s' добавлена новая глава!")
	}

	return nil
}

func (s *ModerationService) RejectChapter(ctx context.Context, id string) error {
	return s.NovelaRepo.UpdateChapterStatus(ctx, id, "rejected")
}

func (s *ModerationService) UpdateVolume(ctx context.Context, id string, req dto.UpdateVolumeRequest) error {
	return s.NovelaRepo.UpdateVolume(ctx, id, req.VolumeNumber, req.Title)
}

func (s *ModerationService) UpdateChapter(ctx context.Context, id string, req dto.UpdateChapterRequest) error {
	return s.NovelaRepo.UpdateChapter(ctx, id, req.ChapterNumber, req.Title, req.Content)
}
