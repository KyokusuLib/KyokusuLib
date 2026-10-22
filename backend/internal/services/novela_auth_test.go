package service

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/repository"
)

func newAuthTestService(t *testing.T) *NovelaService {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return NewNovelaService(
		repository.NewNovelaRepository(db),
		repository.NewNovelaRatingRepository(db),
		repository.NewNovelaBookmarkRepository(db),
		repository.NewNovelaLikeRepository(db),
		repository.NewUserRepository(db),
		repository.NewUserProfileRepository(db),
		repository.NewUserTagRepository(db),
		nil,
		nil,
	)
}

func unprivilegedUserID(t *testing.T, s *NovelaService, novelaID int) int {
	t.Helper()

	const query = `
		SELECT u.id
		FROM users u
		WHERE u.role NOT IN ('admin', 'moderator')
		  AND NOT EXISTS (
			SELECT 1
			FROM novela_teams nt
			JOIN team_members tm ON tm.team_id = nt.team_id
			WHERE nt.novela_id = $1
			  AND tm.user_id = u.id
			  AND tm.role IN ('owner', 'admin', 'moderator')
		  )
		LIMIT 1`

	var id int
	if err := s.Repo.DB.QueryRow(query, novelaID).Scan(&id); err != nil {
		t.Skipf("no unprivileged user in test database: %v", err)
	}
	return id
}

func adminUserID(t *testing.T, s *NovelaService) int {
	t.Helper()

	var id int
	if err := s.Repo.DB.QueryRow(`SELECT id FROM users WHERE role = 'admin' LIMIT 1`).Scan(&id); err != nil {
		t.Skipf("no admin user in test database: %v", err)
	}
	return id
}

func TestUpdateNovelaRejectsUserWithoutPermission(t *testing.T) {
	s := newAuthTestService(t)

	const novelaID = 7
	userID := unprivilegedUserID(t, s, novelaID)

	err := s.UpdateNovela(context.Background(), novelaID, userID, dto.UpdateNovelaRequest{}, "")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUpdateNovelaAllowsAdmin(t *testing.T) {
	s := newAuthTestService(t)
	userID := adminUserID(t, s)

	err := s.UpdateNovela(context.Background(), 999999, userID, dto.UpdateNovelaRequest{}, "")
	if !errors.Is(err, ErrNovelaNotFound) {
		t.Fatalf("expected ErrNovelaNotFound, got %v", err)
	}
}

func TestUpdateNovelaUnknownUserIsForbidden(t *testing.T) {
	s := newAuthTestService(t)

	err := s.UpdateNovela(context.Background(), 7, 999999, dto.UpdateNovelaRequest{}, "")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
