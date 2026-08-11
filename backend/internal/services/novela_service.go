package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/lanxre/kyokusulib/internal/constants"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"github.com/lanxre/kyokusulib/internal/strategy"
	"github.com/lanxre/kyokusulib/internal/utils/files"
	"github.com/redis/go-redis/v9"
)

type NovelaService struct {
	Repo                *repository.NovelaRepository
	RatingRepo          *repository.NovelaRatingRepository
	BookmarkRepo        *repository.NovelaBookmarkRepository
	LikeRepo            *repository.NovelaLikeRepository
	UserRepo            *repository.UserRepository
	UserProfileRepo     *repository.UserProfileRepository
	UserTagRepo         *repository.UserTagRepository
	NotificationService *NotificationService
	redis          		*redis.Client
	awardStrategies     []strategy.AwardStrategy
}

func NewNovelaService(repo *repository.NovelaRepository,
	novelaRatingRepo *repository.NovelaRatingRepository,
	novelaBookmarkRepo *repository.NovelaBookmarkRepository,
	novelaLikeRepo *repository.NovelaLikeRepository,
	userRepo *repository.UserRepository,
	userProfileRepo *repository.UserProfileRepository,
	userTagRepo *repository.UserTagRepository,
	notificationService *NotificationService,
	redisClient *redis.Client,
) *NovelaService {

	ctx := context.Background()
	
	return &NovelaService{
		Repo:                repo,
		RatingRepo:          novelaRatingRepo,
		BookmarkRepo:        novelaBookmarkRepo,
		LikeRepo:            novelaLikeRepo,
		UserRepo:            userRepo,
		UserProfileRepo:     userProfileRepo,
		NotificationService: notificationService,
		redis:               redisClient,
		awardStrategies:     initStrategies(ctx, userTagRepo, userRepo),
	}
}

func initStrategies(
	ctx context.Context,
	userTagRepo *repository.UserTagRepository,
	userRepo *repository.UserRepository,
) []strategy.AwardStrategy {

	strategies := []strategy.AwardStrategy{}

	// при прочтение десяти глав
	if tag, err := userTagRepo.GetUserTagById(ctx, int(constants.ReadTenChaptersTagId)); err == nil {
		strategies = append(strategies, strategy.NewTagAward(
			string(constants.ReadTenChaptersTagName),
			tag.ID,
			10,
			strategy.EventReadChapter,
			userRepo,
		))
	}

	// при получение второго уровня
	if tag, err := userTagRepo.GetUserTagById(ctx, int(constants.LevelTwoNewbieTagId)); err == nil {
		strategies = append(strategies, strategy.NewTagAward(
			string(constants.LevelTwoNewbieTagName),
			tag.ID,
			2,
			strategy.EventLevelUp,
			userRepo,
		))
	}

	return strategies
}

func (s *NovelaService) GetNovelaById(ctx context.Context, id, userID int) (dto.NovelaResponse, error) {
	tx, err := s.Repo.DB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return dto.NovelaResponse{}	, err
	}
	defer tx.Rollback()

	novela, err := s.Repo.GetFullByID(tx, ctx, id, userID)
	if err != nil || novela == nil {
		return dto.NovelaResponse{}, ErrNovelaNotFound
	}

	ratingData, _ := s.RatingRepo.GetRating(tx, ctx, novela.ID)
	bookmarkData, _ := s.BookmarkRepo.GetBookmarkStats(tx, ctx, novela.ID)

	res := s.novelaToDto(novela)
	res.RatingDetails = s.mapRatingToDto(ratingData)
	res.BookmarkDetails = s.mapBookmarkToDto(bookmarkData)

	return res, nil
}

func (s *NovelaService) UploadPoster(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	return files.UploadImage(ctx, file, header, "novelas/posters", 600, 900)
}

func (s *NovelaService) SavePoster(ctx context.Context, id int, posterURL string) error {
	return s.Repo.UpdatePoster(ctx, id, posterURL)
}

func (s *NovelaService) Create(ctx context.Context, novela *db.Novela) error {
	return s.Repo.Create(ctx, novela)
}

func (s *NovelaService) SetBookmark(ctx context.Context, userID int, req dto.UpdateBookmarkRequest) error {
	categoryID := req.CategoryID
	if categoryID == 0 && req.Category != "" {
		id, err := s.BookmarkRepo.GetCategoryByName(ctx, userID, req.Category)
		if err != nil {
			return errors.New("invalid category")
		}
		categoryID = id
	} else if categoryID == 0 {
		return errors.New("category is required")
	}

	return s.BookmarkRepo.SetBookmark(ctx, userID, req.NovelaID, categoryID)
}

func (s *NovelaService) RemoveBookmark(ctx context.Context, userID, novelaID int) error {
	return s.BookmarkRepo.RemoveBookmark(ctx, userID, novelaID)
}

func (s *NovelaService) SetLike(ctx context.Context, userID int, req dto.UpdateLikeRequest) error {
	return s.LikeRepo.SetLike(ctx, &db.NovelaLike{
		UserID:   userID,
		NovelaID: req.NovelaID,
		HasLiked: req.HasLiked,
	})
}

func (s *NovelaService) SetRating(ctx context.Context, userID int, req dto.UpdateRatingRequest) error {
	return s.RatingRepo.SetRating(ctx, &db.NovelaRating{
		UserID:   userID,
		NovelaID: req.NovelaID,
		Rating:   req.Rating,
	})
}

func (s *NovelaService) UpdateNovela(ctx context.Context, id int, req dto.UpdateNovelaRequest, posterURL string) error {
	releaseDate, _ := time.Parse("2006", req.ReleaseDate)
	novela := &db.Novela{
		ID:                id,
		Title:             req.Title,
		AlternativeTitles: req.AlternativeTitles,
		Description:       req.Description,
		Type:              req.Type,
		AgeRating:         req.AgeRating,
		ReleaseDate:       releaseDate,
		Status:            req.Status,
		Country:           req.Country,
		TranslationStatus: req.TranslationStatus,
		PosterURL:         posterURL,
		Genres:            req.Genres,
		Categories:        req.Categories,
	}

	for _, authorID := range req.Authors {
		novela.Authors = append(novela.Authors, db.NovelaAuthor{ID: authorID})
	}

	err := s.Repo.UpdateFull(ctx, novela)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.ErrUnsupported
		}
		return ErrInternal
	}

	return nil
}

func (s *NovelaService) GetNovelas(ctx context.Context, userID int, f dto.NovelaFilters) ([]dto.NovelaResponse, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}

	tx, err := s.Repo.DB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	dbNovelas, total, err := s.Repo.GetNovelas(tx, ctx, userID, f)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.NovelaResponse, 0, len(dbNovelas))
	novelaIDs := make([]int, 0, len(dbNovelas))
	for _, n := range dbNovelas {
		novelaIDs = append(novelaIDs, n.ID)
	}

	ratings, err := s.RatingRepo.GetRatingsBatch(tx, ctx, novelaIDs)
	if err != nil {
		return nil, 0, err
	}
	bookmarks, err := s.BookmarkRepo.GetBookmarkStatsBatch(tx, ctx, novelaIDs)
	if err != nil {
		return nil, 0, err
	}

	for _, n := range dbNovelas {
		d := s.novelaToDto(&n)
		d.RatingDetails = s.mapRatingToDto(ratings[d.ID])
		d.BookmarkDetails = s.mapBookmarkToDto(bookmarks[d.ID])

		res = append(res, d)
	}

	return res, total, tx.Commit()
}

func (s *NovelaService) GetUserNovelaBookmarks(ctx context.Context, userID int, category string) ([]dto.UserNovelaBookmark, error) {
	var categoryID int
	var err error

	// If the category is a number, it's already an ID
	if id, parseErr := strconv.Atoi(category); parseErr == nil {
		categoryID = id
	} else {
		// Otherwise, look it up by name
		categoryID, err = s.BookmarkRepo.GetCategoryByName(ctx, userID, category)
		if err != nil {
			return nil, errors.New("Invalid category type")
		}
	}

	novelas, err := s.Repo.GetUserNovelaBookmarks(ctx, userID, categoryID)
	if err != nil {
		return nil, errors.New("Invalid category type ")
	}

	var ns []dto.UserNovelaBookmark
	for _, n := range novelas {
		ns = append(ns, dto.UserNovelaBookmark{
			ID:        n.ID,
			Title:     n.Title,
			PosterURL: n.PosterURL,
			Type:      n.Type,
			Rating:    n.Rating,
		})
	}

	return ns, nil
}

func (s *NovelaService) novelaToDto(n *db.Novela) dto.NovelaResponse {
	authors := make([]dto.NovelaAuthor, 0)
	for _, a := range n.Authors {
		authors = append(authors, dto.NovelaAuthor{ID: a.ID, Name: a.Name, Role: a.Role})
	}

	volumes := make([]dto.NovelaVolume, 0)
	for _, v := range n.Volumes {
		chapters := make([]dto.NovelaChapter, 0)
		for _, ch := range v.Chapters {
			imgs := make([]dto.NovelaChapterImage, 0)
			for _, img := range ch.Images {
				imgs = append(imgs, dto.NovelaChapterImage{ID: img.ID, ImageURL: img.ImageURL, Caption: img.Caption, Position: img.Position})
			}
			chapters = append(chapters, dto.NovelaChapter{ID: ch.ID, Title: ch.Title, Number: ch.Number, Images: imgs, IsRead: ch.IsRead})
		}
		volumes = append(volumes, dto.NovelaVolume{ID: v.ID, Title: v.Title, Number: v.Number, Chapters: chapters})
	}

	var bookmark *string
	if n.Bookmark != nil {
		str := n.Bookmark.Name
		bookmark = &str
	}

	var lastReadChapterID *string
	var lastReadChapterNumber *float64
	var lastReadChapterScroll *int
	if n.LastReadChapter != nil {
		lastReadChapterID = &n.LastReadChapter.ID
		lastReadChapterNumber = &n.LastReadChapter.Number
		lastReadChapterScroll = &n.LastReadChapter.ScrollPosition
	}

	return dto.NovelaResponse{
		ID: n.ID, Title: n.Title, AlternativeTitles: n.AlternativeTitles, Description: n.Description,
		Type: n.Type, AgeRating: n.AgeRating, ReleaseDate: n.ReleaseDate.Format("2006-01-02"),
		Status: n.Status, TranslationStatus: n.TranslationStatus, PosterURL: n.PosterURL,
		Country: n.Country, Views: n.Views, Genres: n.Genres, Categories: n.Categories,
		Authors: authors, Volumes: volumes, Bookmark: bookmark, HasLiked: n.HasLiked,
		LikeCount: n.LikeCount, UserRating: n.UserRating,
		LastReaded: &dto.NovelaLastReaded{ChapterID: lastReadChapterID, ChapterNumber: lastReadChapterNumber, ChapterScroll: lastReadChapterScroll},
	}
}

func (s *NovelaService) mapRatingToDto(data *db.NovelaRatingSummary) dto.NovelaRatingCategory {
	items := make([]dto.NCItem, 0)
	if data == nil {
		return dto.NovelaRatingCategory{NCItems: items}
	}

	for i := 10; i >= 1; i-- {
		count := 0
		if val, ok := data.Distribution[i]; ok {
			count = val
		}
		items = append(items, dto.NCItem{Value: i, Count: count})
	}
	return dto.NovelaRatingCategory{Total: data.TotalCount, TotalRating: data.AverageRating, NCItems: items}
}

func (s *NovelaService) mapBookmarkToDto(data *db.NovelaBookmarkSummary) dto.NovelaBookmarkCategory {
	items := make([]dto.NCItem, 0)
	if data == nil {
		return dto.NovelaBookmarkCategory{NCItems: items}
	}

	for catName, count := range data.Distribution {
		items = append(items, dto.NCItem{Value: catName, Count: count})
	}
	return dto.NovelaBookmarkCategory{Total: data.TotalCount, NCItems: items}
}

func (s *NovelaService) GetBookmarkCategories(ctx context.Context, userID int) ([]db.BookmarkCategoryCount, error) {
	return s.BookmarkRepo.GetCategoriesWithCount(ctx, userID)
}

func (s *NovelaService) CreateBookmarkCategory(ctx context.Context, userID int, name string) (int, error) {
	return s.BookmarkRepo.CreateCategory(ctx, userID, name)
}

func (s *NovelaService) UpdateBookmarkCategory(ctx context.Context, id int, userID int, name string, visible bool) error {
	return s.BookmarkRepo.UpdateCategory(ctx, id, userID, name, visible)
}

func (s *NovelaService) DeleteBookmarkCategory(ctx context.Context, id int, userID int) error {
	return s.BookmarkRepo.DeleteCategory(ctx, id, userID)
}

func (s *NovelaService) DeleteNovela(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}

func (s *NovelaService) AddTeamToNovela(ctx context.Context, novelaID, teamID int) error {
	count, err := s.Repo.CountNovelaTeams(ctx, novelaID)
	if err != nil {
		return err
	}
	if count >= 5 {
		return errors.New("maximum 5 teams can be added to a novela")
	}
	return s.Repo.AddTeamToNovela(ctx, novelaID, teamID)
}

func (s *NovelaService) AddVolume(ctx context.Context, novelaID int, userID int, req dto.AddVolumeRequest) (string, string, error) {
	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return "", "", err
	}

	hasTeamPermission, err := s.Repo.CheckUserNovelaTeamPermission(ctx, userID, novelaID)
	if err != nil {
		return "", "", err
	}

	status := "pending"
	if user.Role == "admin" || user.Role == "moderator" || hasTeamPermission {
		status = "approved"
	}

	id, err := s.Repo.AddVolume(ctx, novelaID, req.VolumeNumber, req.Title, status, userID)
	if err == nil && status == "approved" {
		go s.safeSendNovelaNotification(novelaID, "Обновление", "В новеллу '%s' добавлен новый том!")
	}
	return id, status, err
}

func (s *NovelaService) AddChapter(ctx context.Context, volumeID string, userID int, req dto.AddChapterRequest) (string, string, error) {
	novelaID, err := s.Repo.GetNovelaIDByVolumeID(ctx, volumeID)
	if err != nil {
		return "", "", err
	}

	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return "", "", err
	}

	hasTeamPermission, err := s.Repo.CheckUserNovelaTeamPermission(ctx, userID, novelaID)
	if err != nil {
		return "", "", err
	}

	status := "pending"
	if user.Role == "admin" || user.Role == "moderator" || hasTeamPermission {
		status = "approved"
	}

	id, err := s.Repo.AddChapter(ctx, volumeID, req.ChapterNumber, req.Title, req.Content, status, userID)
	if err == nil && status == "approved" {
		go s.safeSendNovelaNotification(novelaID, "Обновление", "В новеллу '%s' добавлена новая глава!")
	}
	return id, status, err
}

func (s *NovelaService) sendNovelaNotification(novelaID int, title, messageFmt string) {
	bgCtx := context.Background()
	novelaTitle, err := s.Repo.GetTitleByID(bgCtx, novelaID)
	if err != nil {
		return
	}
	userIDs, err := s.BookmarkRepo.GetUsersWithNovelaBookmark(bgCtx, novelaID)
	if err != nil {
		return
	}
	message := fmt.Sprintf(messageFmt, novelaTitle)

	intIDs := make([]int64, len(userIDs))
	for i, userID := range userIDs {
		intIDs[i] = int64(userID)
	}

	if err := s.NotificationService.CreateBatch(bgCtx, intIDs, title, message); err != nil {
		return
	}
}

func (s *NovelaService) safeSendNovelaNotification(novelaID int, title, messageFmt string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in sendNovelaNotification: %v", r)
		}
	}()
	s.sendNovelaNotification(novelaID, title, messageFmt)
}

func (s *NovelaService) DeleteChapterImages(ctx context.Context, chapterID string) error {
	return s.Repo.DeleteChapterImages(ctx, chapterID)
}

func (s *NovelaService) AddChapterImage(ctx context.Context, chapterID string, req dto.AddChapterImageRequest) (int, error) {
	return s.Repo.AddChapterImage(ctx, chapterID, req.ImageURL, req.Caption, req.Position)
}

func (s *NovelaService) AddChapterImages(ctx context.Context, chapterID string, images []dto.AddChapterImageRequest) error {
	return s.Repo.AddChapterImages(ctx, chapterID, images)
}

func (s *NovelaService) GetChapterReaderDetails(ctx context.Context, chapterID string, userID int) (*dto.ChapterReaderResponse, error) {
	return s.Repo.GetChapterReaderDetails(ctx, chapterID, userID)
}

func (s *NovelaService) SaveChapterReadPosition(ctx context.Context, userID int, req dto.SaveReadPositionRequest) error {
	return s.Repo.UpdateChapterReadPosition(ctx, userID, req.ChapterID, req.ScrollPosition)
}

func (s *NovelaService) UpdateChapter(ctx context.Context, chapterID string, userID int, req dto.AddChapterRequest) error {
	canManage, err := s.canManageChapter(ctx, userID, chapterID)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("forbidden: you do not have permission to update this chapter")
	}

	return s.Repo.UpdateChapter(ctx, chapterID, req.ChapterNumber, req.Title, req.Content)
}

func (s *NovelaService) DeleteVolume(ctx context.Context, volumeID string, userID int) error {
	novelaID, err := s.Repo.GetNovelaIDByVolumeID(ctx, volumeID)
	if err != nil {
		return err
	}

	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return err
	}

	hasTeamPermission, err := s.Repo.CheckUserNovelaTeamPermission(ctx, userID, novelaID)
	if err != nil {
		return err
	}

	if user.Role != "admin" && user.Role != "moderator" && !hasTeamPermission {
		return errors.New("forbidden: you do not have permission to delete this volume")
	}

	return s.Repo.DeleteVolume(ctx, volumeID)
}

func (s *NovelaService) DeleteChapter(ctx context.Context, chapterID string, userID int) error {
	canManage, err := s.canManageChapter(ctx, userID, chapterID)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("forbidden: you do not have permission to delete this chapter")
	}

	return s.Repo.DeleteChapter(ctx, chapterID)
}

func (s *NovelaService) MarkChapterAsRead(ctx context.Context, userID int, chapterID string) (*db.UserLevel, error) {
	exist, err := s.Repo.IsExistChapter(ctx, chapterID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("Chapter not found")
	}

	read, err := s.Repo.IsReadChapter(ctx, userID, chapterID)
	if err != nil || read {
		return nil, errors.New("Chapter already read")
	}

	if userID > 0 {
		if _, err := s.Repo.MarkChapterAsRead(ctx, userID, chapterID); err != nil {
			return nil, err
		}
	}
	
	dataLevel, err := s.UpLevelForReading(ctx, userID)
	if err := s.AwardUser(ctx, userID, strategy.EventReadChapter, nil); err != nil {
		return nil, err
	}
	return dataLevel, nil
}

func (s *NovelaService) canManageChapter(ctx context.Context, userID int, chapterID string) (bool, error) {
	chapter, err := s.Repo.GetChapterByID(ctx, chapterID)
	if err != nil {
		return false, err
	}
	if chapter == nil {
		return false, errors.New("chapter not found")
	}

	novelaID, err := s.Repo.GetNovelaIDByVolumeID(ctx, chapter.VolumeID)
	if err != nil {
		return false, err
	}

	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	if user.Role == "admin" || user.Role == "moderator" {
		return true, nil
	}

	hasTeamPermission, err := s.Repo.CheckUserNovelaTeamPermission(ctx, userID, novelaID)
	if err != nil {
		return false, err
	}

	return hasTeamPermission, nil
}

func (s *NovelaService) GetMostSearched(ctx context.Context, limit int) (*dto.MostSearchedResponse, error) {
	cacheKey := fmt.Sprintf("most_searched:%d", limit)

	if s.redis != nil {
		val, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedResponse dto.MostSearchedResponse
			if err := json.Unmarshal([]byte(val), &cachedResponse); err == nil {
				return &cachedResponse, nil
			}
		}
	}
	
	genres, err := s.Repo.GetMostSearchedGenres(ctx, limit)
	if err != nil {
		return nil, err
	}

	categories, err := s.Repo.GetMostSearchedCategories(ctx, limit)
	if err != nil {
		return nil, err
	}

	res := &dto.MostSearchedResponse{
		Genres: genres,
		Categories: categories,
	}

	if s.redis != nil {
		if data, err := json.Marshal(res); err == nil {
			_ = s.redis.Set(ctx, cacheKey, data, 15 * time.Minute).Err()
		}
	}

	return res, nil
}

func (s *NovelaService) UpLevelForReading(ctx context.Context, userID int) (*db.UserLevel, error) {
	userProfile, err := s.UserProfileRepo.GetUserLevel(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("User not found. cannot uplevel: %w", err)
	}

	var experience int64
	leveledUp := false
	

	totalLevel, err := s.UserProfileRepo.GetLevelDefinition(ctx, userProfile.Level + 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get level definition: %w", err)
	}
	
	if totalLevel == nil {
		return nil, fmt.Errorf("level definition not found for level %d", userProfile.Level + 1)
	}

	if constants.UpLevelExperienceThreshold >= userProfile.XPForNext {
		s.UserProfileRepo.UpUserLevel(ctx, userID)
		leveledUp = true
		userProfile.LevelTitle = totalLevel.Title
		userProfile.Level = totalLevel.Level
		experience = (userProfile.Experience + constants.UpLevelExperienceThreshold) - int64(totalLevel.TotalXpRequired)
	} else {
		experience = userProfile.Experience + constants.UpLevelExperienceThreshold
	}

	s.UserProfileRepo.SetUserExperiance(ctx, userID, experience)

	if leveledUp {
		s.AwardUser(ctx, userID, strategy.EventLevelUp, map[string]any{
			"level": totalLevel.Level,
		})
		s.NotificationService.Create(ctx, int64(userID), "Повышение уровня", fmt.Sprintf("Вы достигли уровня %d - (%s)", userProfile.Level, userProfile.LevelTitle))
	}

	return &db.UserLevel{
		Level:         userProfile.Level,
		LevelTitle:    userProfile.LevelTitle,
		Experience:    experience,
		XPForNext:     int64(totalLevel.TotalXpRequired) - experience,
	}, nil
}

func (s *NovelaService) AwardUser(ctx context.Context, userID int, eventType string, extra map[string]any) error {
	for _, st := range s.awardStrategies {
		if err := st.Evaluate(ctx, strategy.AwardEvent{
			Type:   eventType,
			UserID: userID,
			Extra:  extra,
		}); err != nil {
			return err
		}
	}
	
	return nil
}
