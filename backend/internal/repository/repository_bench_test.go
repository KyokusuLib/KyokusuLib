package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lanxre/kyokusulib/internal/models/dto"
)

const benchNovelaID = 7

func benchDB(b *testing.B) *sql.DB {
	b.Helper()

	url := os.Getenv("BENCH_DATABASE_URL")
	if url == "" {
		b.Skip("BENCH_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.Ping(); err != nil {
		b.Fatal(err)
	}

	b.Cleanup(func() { _ = db.Close() })
	return db
}

func BenchmarkNovelaGetFullByID(b *testing.B) {
	db := benchDB(b)
	repo := NewNovelaRepository(db)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
		if err != nil {
			b.Fatal(err)
		}
		if _, err := repo.GetFullByID(tx, ctx, benchNovelaID, 0); err != nil {
			_ = tx.Rollback()
			b.Fatal(err)
		}
		_ = tx.Rollback()
	}
}

func BenchmarkNovelaPage(b *testing.B) {
	db := benchDB(b)
	novelaRepo := NewNovelaRepository(db)
	ratingRepo := NewNovelaRatingRepository(db)
	bookmarkRepo := NewNovelaBookmarkRepository(db)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
		if err != nil {
			b.Fatal(err)
		}
		novela, err := novelaRepo.GetFullByID(tx, ctx, benchNovelaID, 0)
		if err != nil {
			_ = tx.Rollback()
			b.Fatal(err)
		}
		if novela != nil {
			_, _ = ratingRepo.GetRating(tx, ctx, novela.ID)
			_, _ = bookmarkRepo.GetBookmarkStats(tx, ctx, novela.ID)
		}
		_ = tx.Rollback()
	}
}

func BenchmarkNovelaList(b *testing.B) {
	db := benchDB(b)
	repo := NewNovelaRepository(db)
	ctx := context.Background()
	filters := dto.NovelaFilters{Limit: 20}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
		if err != nil {
			b.Fatal(err)
		}
		if _, _, err := repo.GetNovelas(tx, ctx, 0, filters); err != nil {
			_ = tx.Rollback()
			b.Fatal(err)
		}
		_ = tx.Rollback()
	}
}

func BenchmarkTgPostList(b *testing.B) {
	db := benchDB(b)
	repo := NewTgPostRepository(db)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.List(ctx, 20, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserGetUsers(b *testing.B) {
	db := benchDB(b)
	repo := NewUserRepository(db)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.GetUsers(ctx, "", 20, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserListServiceBatches(b *testing.B) {
	db := benchDB(b)
	userRepo := NewUserRepository(db)
	profileRepo := NewUserProfileRepository(db)
	ctx := context.Background()

	users, err := userRepo.GetUsers(ctx, "", 20, 0)
	if err != nil {
		b.Fatal(err)
	}
	ids := make([]int, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := userRepo.GetUsers(ctx, "", 20, 0); err != nil {
			b.Fatal(err)
		}
		if _, err := profileRepo.GetUserLevelsBatch(ctx, ids); err != nil {
			b.Fatal(err)
		}
		if _, err := userRepo.GetUserTagsBatch(ctx, ids); err != nil {
			b.Fatal(err)
		}
		if _, err := userRepo.GetUserStatsBatch(ctx, ids); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNotificationList(b *testing.B) {
	db := benchDB(b)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.GetByUserID(ctx, 8, 20, 0); err != nil {
			b.Fatal(err)
		}
	}
}
