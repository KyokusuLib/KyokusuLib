package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/constants"
	"github.com/lanxre/kyokusulib/internal/middleware"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	service "github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/utils/files"
	"github.com/lanxre/kyokusulib/internal/utils/response"
)

type NovelaHandler struct {
	service   *service.NovelaService
	Validator *validator.Validate
}

func NewNovelaHandler(ranobeService *service.NovelaService, validator *validator.Validate) *NovelaHandler {
	return &NovelaHandler{
		service:   ranobeService,
		Validator: validator,
	}
}

func (h *NovelaHandler) GetNovela(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	includeChapters := true
	if summary, err := strconv.ParseBool(r.URL.Query().Get("summary")); err == nil && summary {
		includeChapters = false
	}

	novela, err := h.service.GetNovelaById(r.Context(), id, userID, includeChapters)
	if err != nil {
		if errors.Is(err, service.ErrNovelaNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, novela)
}

func (h *NovelaHandler) GetVolumeChapters(w http.ResponseWriter, r *http.Request) {
	volumeID := mux.Vars(r)["id"]
	if _, err := uuid.Parse(volumeID); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid volume id")
		return
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	chapters, err := h.service.GetVolumeChapters(r.Context(), volumeID, userID)
	if err != nil {
		log.Printf("get volume chapters: %v", err)
		response.Error(w, http.StatusInternalServerError, "Failed to load chapters")
		return
	}
	response.JSON(w, http.StatusOK, chapters)
}

func (h *NovelaHandler) CreateNovela(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "File too big or invalid format")
		return
	}

	input := dto.CreateNovelaDTO{
		Title:             r.FormValue("title"),
		Description:       r.FormValue("description"),
		Type:              r.FormValue("type"),
		AgeRating:         r.FormValue("age_rating"),
		ReleaseYear:       r.FormValue("release_year"),
		Status:            r.FormValue("status"),
		Country:           r.FormValue("country"),
		TranslationStatus: r.FormValue("translation_status"),
	}

	_ = json.Unmarshal([]byte(r.FormValue("genres")), &input.Genres)
	_ = json.Unmarshal([]byte(r.FormValue("categories")), &input.Categories)

	var altTitles []string
	_ = json.Unmarshal([]byte(r.FormValue("alternative_titles")), &altTitles)

	if err := h.Validator.Struct(input); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	var posterURL string
	file, header, err := r.FormFile("poster")
	if err == nil {
		defer file.Close()

		url, err := h.service.UploadPoster(r.Context(), file, header)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		posterURL = url
	}

	releaseDate, _ := time.Parse("2006", input.ReleaseYear)

	novela := &db.Novela{
		Title:             input.Title,
		AlternativeTitles: altTitles,
		Description:       input.Description,
		Type:              input.Type,
		AgeRating:         input.AgeRating,
		ReleaseDate:       releaseDate,
		Status:            input.Status,
		Country:           input.Country,
		TranslationStatus: input.TranslationStatus,
		PosterURL:         posterURL,
		Genres:            input.Genres,
		Categories:        input.Categories,
	}

	if err := h.service.Create(r.Context(), novela); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"id":      novela.ID,
		"message": "Novela created successfully",
	})
}

func (h *NovelaHandler) UpdateNovela(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid novela id")
		return
	}

	if err := h.service.AuthorizeNovelaEdit(r.Context(), id, userID); err != nil {
		h.writeUpdateError(w, id, userID, err)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "Payload error")
		return
	}

	var req dto.UpdateNovelaRequest
	if err := json.Unmarshal([]byte(r.FormValue("data")), &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	var posterURL string
	if file, header, err := r.FormFile("poster"); err == nil {
		defer file.Close()
		posterURL, _ = h.service.UploadPoster(r.Context(), file, header)
	}

	err = h.service.UpdateNovela(r.Context(), id, userID, req, posterURL)
	if err != nil {
		if posterURL != "" {
			_ = files.DeleteImage(posterURL)
		}
		h.writeUpdateError(w, id, userID, err)
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) writeUpdateError(w http.ResponseWriter, id, userID int, err error) {
	switch {
	case errors.Is(err, service.ErrForbidden):
		response.Error(w, http.StatusForbidden, "Forbidden")
	case errors.Is(err, service.ErrNovelaNotFound):
		response.Error(w, http.StatusNotFound, "Novela not found")
	default:
		log.Printf("update novela %d by user %d: %v", id, userID, err)
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *NovelaHandler) SetBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "User not found")
		return
	}

	var req dto.UpdateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ErrorWithDetails(w, http.StatusBadRequest, map[string]string{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.SetBookmark(r.Context(), userID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) RemoveBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "User not found")
		return
	}

	vars := mux.Vars(r)
	novelaIDStr := vars["id"]

	novelaID, err := strconv.Atoi(novelaIDStr)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid novela ID format")
		return
	}

	if novelaID <= 0 {
		response.Error(w, http.StatusBadRequest, "ID must be greater than zero")
		return
	}

	if err := h.service.RemoveBookmark(r.Context(), userID, novelaID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) SetLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "User not found")
		return
	}

	var req dto.UpdateLikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ErrorWithDetails(w, http.StatusBadRequest, map[string]string{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.SetLike(r.Context(), userID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) SetRating(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "User not found")
		return
	}

	var req dto.UpdateRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ErrorWithDetails(w, http.StatusBadRequest, map[string]string{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.SetRating(r.Context(), userID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) GetNovelas(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)
	q := r.URL.Query()

	toInt := func(k string, d int) int {
		v, err := strconv.Atoi(q.Get(k))
		if err != nil {
			return d
		}
		return v
	}

	parseCSV := func(k string) []string {
		trimQuotes := func(s string) string {
			return strings.Trim(s, "'\" ")
		}

		if v := q.Get(k); v != "" {
			parts := strings.Split(v, ",")
			for i := range parts {
				parts[i] = trimQuotes(parts[i])
			}
			return parts
		}

		values := q[k]
		for i := range values {
			values[i] = trimQuotes(values[i])
		}
		return values
	}

	filters := dto.NovelaFilters{
		Limit:             toInt("limit", 20),
		Offset:            toInt("offset", 0),
		Search:            q.Get("search"),
		Sort:              dto.NovelaSort(q.Get("sort")),
		Type:              q.Get("type"),
		Status:            q.Get("status"),
		TranslationStatus: q.Get("translation_status"),
		ChaptersFrom:      toInt("chapters_from", 0),
		ChaptersTo:        toInt("chapters_to", 0),
		YearFrom:          toInt("year_from", 0),
		YearTo:            toInt("year_to", 0),
		Genres:            parseCSV("genres"),
		Categories:        parseCSV("categories"),
		Country:           parseCSV("country"),
		AgeRating:         q.Get("ageRating"),
		AuthorID:          toInt("author_id", 0),
	}

	novelas, total, err := h.service.GetNovelas(r.Context(), userID, filters)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	response.JSON(w, http.StatusOK, novelas)
}

func (h *NovelaHandler) GetUserNovelaBookmarks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userID := query.Get("user_id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "Missing user_id parameter")
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	category := query.Get("category")
	if category == "" {
		response.Error(w, http.StatusBadRequest, "Missing category parameter")
		return
	}

	novelas, err := h.service.GetUserNovelaBookmarks(r.Context(), userIDInt, category)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithEntity(w, http.StatusOK, novelas)
}

func (h *NovelaHandler) GetBookmarkCategories(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := r.URL.Query()
	userID := query.Get("user_id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "Missing user_id parameter")
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	categories, err := h.service.GetBookmarkCategories(r.Context(), userIDInt)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithEntity(w, http.StatusOK, categories)
}

func (h *NovelaHandler) CreateBookmarkCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.BookmarkCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	id, err := h.service.CreateBookmarkCategory(r.Context(), userID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"id":      id,
		"message": "Category created successfully",
	})
}

func (h *NovelaHandler) UpdateBookmarkCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	categoryIDStr := vars["id"]
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req dto.BookmarkCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	visible := true
	if req.Visible != nil {
		visible = *req.Visible
	}

	if err := h.service.UpdateBookmarkCategory(r.Context(), categoryID, userID, req.Name, visible); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) DeleteBookmarkCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	categoryIDStr := vars["id"]
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.DeleteBookmarkCategory(r.Context(), categoryID, userID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) AddTeamToNovela(w http.ResponseWriter, r *http.Request) {
	novelaID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.AddNovelaTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	if err := h.service.AddTeamToNovela(r.Context(), novelaID, req.TeamID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) AddVolume(w http.ResponseWriter, r *http.Request) {
	novelaID, _ := strconv.Atoi(mux.Vars(r)["id"])
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	var req dto.AddVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	id, status, err := h.service.AddVolume(r.Context(), novelaID, userID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := "Successfully added"
	if status == "pending" {
		msg = "Sent for moderation"
	}

	response.JSON(w, http.StatusCreated, map[string]any{"id": id, "message": msg, "status": status})
}

func (h *NovelaHandler) AddChapter(w http.ResponseWriter, r *http.Request) {
	volumeID := mux.Vars(r)["id"]
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	var req dto.AddChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	id, status, err := h.service.AddChapter(r.Context(), volumeID, userID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := "Successfully added"
	if status == "pending" {
		msg = "Sent for moderation"
	}

	response.JSON(w, http.StatusCreated, map[string]any{"id": id, "message": msg, "status": status})
}

func (h *NovelaHandler) DeleteChapterImages(w http.ResponseWriter, r *http.Request) {
	chapterID := mux.Vars(r)["id"]
	if err := h.service.DeleteChapterImages(r.Context(), chapterID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Images deleted"})
}

func (h *NovelaHandler) AddChapterImage(w http.ResponseWriter, r *http.Request) {
	chapterID := mux.Vars(r)["id"]
	var req dto.AddChapterImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	id, err := h.service.AddChapterImage(r.Context(), chapterID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"id": id, "message": "Image added"})
}

func (h *NovelaHandler) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	chapterID := mux.Vars(r)["id"]
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	var req dto.AddChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	if err := h.service.UpdateChapter(r.Context(), chapterID, userID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Глава успешно обновлена"})
}

func (h *NovelaHandler) DeleteVolume(w http.ResponseWriter, r *http.Request) {
	volumeID := mux.Vars(r)["id"]
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	if err := h.service.DeleteVolume(r.Context(), volumeID, userID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Том успешно удален"})
}

func (h *NovelaHandler) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	chapterID := mux.Vars(r)["id"]
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	if err := h.service.DeleteChapter(r.Context(), chapterID, userID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Глава успешно удалена"})
}

func (h *NovelaHandler) DeleteNovela(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := h.service.DeleteNovela(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Новела успешно удалена",
	})
}

func (h *NovelaHandler) SaveChapterReadPosition(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.SaveReadPositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	if err := h.service.SaveChapterReadPosition(r.Context(), userID, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessOkEmpty(w)
}

func (h *NovelaHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	res, err := h.service.GetChapterReaderDetails(r.Context(), id, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if res == nil {
		response.Error(w, http.StatusNotFound, "Chapter not found")
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *NovelaHandler) MarkChapterAsRead(w http.ResponseWriter, r *http.Request) {
	chapter_id := mux.Vars(r)["id"]
	userID, _ := r.Context().Value(middleware.UserIDKey).(int)

	dataLevel, err := h.service.MarkChapterAsRead(r.Context(), userID, chapter_id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка: %s", err.Error()))
		return
	}

	response.JSON(w, http.StatusOK, dataLevel)
}

func (h *NovelaHandler) GetMostSearched(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetMostSearched(r.Context(), constants.MostSearchedCategoriesLimit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Невозможно получить данные")
		return
	}
	response.JSON(w, http.StatusOK, res)
}
