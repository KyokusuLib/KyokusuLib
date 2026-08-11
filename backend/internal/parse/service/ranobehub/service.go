package ranobehub

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	rhModels "github.com/lanxre/kyokusulib/internal/parse/models/ranobehub"
	service "github.com/lanxre/kyokusulib/internal/services"
	"golang.org/x/sync/errgroup"
)

type RanobeHubParseService struct {
	NovelaService       *service.NovelaService
}

func NewRanobeHubParseService(novelaService *service.NovelaService) *RanobeHubParseService {
	return &RanobeHubParseService{
		NovelaService:       novelaService,
	}
}


func (s *RanobeHubParseService) Parse(ctx context.Context, rhNovela *rhModels.RanobeHubNovela, userID int) error {
	alternativeTitles := []string{rhNovela.Names.Eng}
	releaseDate := time.Date(rhNovela.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	status := ""
	if rhNovela.Status.Title == "Завершено" {
		status = "completed"
	}

	novela := &db.Novela{
		Title:             rhNovela.Names.Rus,
		AlternativeTitles: alternativeTitles,
		Description:       rhNovela.Description,
		Type:              "Ранобе",
		AgeRating:         "16+",
		ReleaseDate:       releaseDate,
		Status:            status,
		TranslationStatus: rhNovela.Status.Title,
		PosterURL:         rhNovela.Posters.Big,
		Country:           "Япония",
	}

	err := s.NovelaService.Create(ctx, novela)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			existingID, lookupErr := s.NovelaService.Repo.GetIDByTitle(ctx, rhNovela.Names.Rus)
			if lookupErr != nil {
				return fmt.Errorf("failed to create novela and failed to find existing one: %w", lookupErr)
			}
			novela.ID = existingID
		} else {
			return fmt.Errorf("failed to create novela: %w", err)
		}
	}

	for _, author := range rhNovela.Authors {
		name := author.Name
		if name == "" {
			name = author.NameEng
		}
		err = s.NovelaService.Repo.LinkAuthor(ctx, novela.ID, name, "Author")
		if err != nil {
			slog.Error("failed to link author", "name", name, "error", err)
		}
	}

	g := new(errgroup.Group)
	for _, vol := range rhNovela.Volumes {
		g.Go(func() error {
			return s.importVolume(ctx, novela.ID, userID, vol)
		})
	}
	return g.Wait()
}

func (s *RanobeHubParseService) importVolume(ctx context.Context, novelaID, userID int, vol rhModels.RanobeHubNovelaVolume) error {
	volumeNumber := int(vol.Num)
	volID, err := s.NovelaService.Repo.GetVolumeIDByNumber(ctx, novelaID, volumeNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			volID, _, err = s.NovelaService.AddVolume(ctx, novelaID, userID, dto.AddVolumeRequest{
				VolumeNumber: volumeNumber,
				Title:        vol.Name,
			})
			if err != nil {
				return fmt.Errorf("failed to add volume %f: %w", vol.Num, err)
			}
		} else {
			return fmt.Errorf("failed to check volume existence: %w", err)
		}
	}

	for _, ch := range vol.Chapters {
		chapterID, err := s.NovelaService.Repo.GetChapterIDByNumber(ctx, volID, ch.Num)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				chapterID, _, err = s.NovelaService.AddChapter(ctx, volID, userID, dto.AddChapterRequest{
					ChapterNumber: ch.Num,
					Title:         ch.Name,
					Content:       ch.Text,
				})
				if err != nil {
					return fmt.Errorf("failed to add chapter %f in volume %f: %w", ch.Num, vol.Num, err)
				}

				images := make([]dto.AddChapterImageRequest, len(ch.Images))
				for imgIdx, imgURL := range ch.Images {
					images[imgIdx] = dto.AddChapterImageRequest{
						ImageURL: imgURL,
						Position: imgIdx,
					}
				}
				if err := s.NovelaService.AddChapterImages(ctx, chapterID, images); err != nil {
					slog.Error("failed to add chapter images", "chapter_id", chapterID, "error", err)
				}
			} else {
				return fmt.Errorf("failed to check chapter existence: %w", err)
			}
		}
	}
	return nil
}

