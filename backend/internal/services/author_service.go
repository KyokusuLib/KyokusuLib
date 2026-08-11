package service

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/lanxre/kyokusulib/internal/constants"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
	"github.com/lanxre/kyokusulib/internal/utils/files"
)

type AuthorService struct {
	repo repository.AuthorRepository
}

func NewAuthorService(repo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{repo: *repo}
}

func (s *AuthorService) CreateAuthor(ctx context.Context, author *db.Author) error {
	return s.repo.CreateAuthor(ctx, author)
}

func (s *AuthorService) UploadAuthorImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	avatarURL, err := files.UploadImage(ctx, file, header, "authors", 400, 400)
    if err != nil {
        return "", err
    }
    
	return avatarURL, nil
}

func (s *AuthorService) GetAuthorByName(ctx context.Context, name string) (*db.Author, error) {
	name = strings.TrimSpace(name)
	name = constants.AuthorNameCleanRegex.ReplaceAllString(name, "")
	name = constants.AuthorNameSpaceRegex.ReplaceAllString(name, " ")
	return s.repo.GetAuthorByName(ctx, name)
}

func (s *AuthorService) GetAuthors(ctx context.Context, search string, limit int, offset int) ([]*db.Author, error) {
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	
	authorsDb, err := s.repo.GetAuthors(ctx, search, limit, offset)

	if err != nil {
		return nil, err
	}

	return authorsDb, nil
}

func (s *AuthorService) GetAuthorById(ctx context.Context, id string) (*dto.AuthorDTO, error) {
	author, err := s.repo.GetAuthorById(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toDTO(author), nil
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, id string, author *db.Author) error {
	return s.repo.UpdateAuthor(ctx, id, author)
}

func (s *AuthorService) toDTO(author *db.Author) *dto.AuthorDTO {
	return &dto.AuthorDTO{
		ID:       author.ID,
		Name:     author.Name,
		Country:  author.Country,
		Metier:   author.Metier,
		Picture:  author.Picture,
		Bio:      author.Bio,
	}
}
