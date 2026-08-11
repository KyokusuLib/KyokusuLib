package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lanxre/kyokusulib/internal/config" 
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

const (
	DefaultTimeout = 5 * time.Second
)

func NewPostgresDB(lc fx.Lifecycle, cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if err := runSeeds(db); err != nil {
		return nil, fmt.Errorf("seeding failed: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection...")
			return db.Close()
		},
	})

	return db, nil
}

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	log.Println("Database migrations check completed")
	return nil
}

func runSeeds(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	log.Println("Checking seed data...")
	
	tagsQuery := `
		INSERT INTO user_tags (tag)
		VALUES 
			('Читатель'),
			('Первые шаги'),
			('Продвинутый читатель'),
			('Читатель-любитель')			
		ON CONFLICT (tag) DO NOTHING;
	`
	if _, err := db.ExecContext(ctx, tagsQuery); err != nil {
		log.Printf("Warning: failed to seed tags: %v", err)
	}

	levelDefinitionsQuery := `
		INSERT INTO level_definitions (level, title, total_xp_required)
		VALUES
			(1, 'Новичок', 0),
			(2, 'Ученик', 1000),
			(3, 'Читатель', 1500),
			(4, 'Практик', 2000),
			(5, 'Эксперт', 2500),
			(6, 'Мастер', 3000),
			(7, 'Глава', 3500),
			(8, 'Величие', 4000)
		ON CONFLICT (level) DO NOTHING;
	`

	if _, err := db.ExecContext(ctx, levelDefinitionsQuery); err != nil {
		log.Printf("Warning: failed to seed level definitions: %v", err)
	}

	genresQuery := `
		INSERT INTO genres (name)
		VALUES 
			('Фэнтези'),
			('Романтика'),
			('Драма'),
			('Приключения'),
			('Детектив'),
			('Фантастика'),
			('Боевик'),
			('Триллер'),
			('Ужасы'),
			('Комедия'),
			('Мистика'),
			('Психология')
		ON CONFLICT (name) DO NOTHING;
	`
	if _, err := db.ExecContext(ctx, genresQuery); err != nil {
		log.Printf("Warning: failed to seed genres: %v", err)
	}

	categoriesQuery := `
		INSERT INTO categories (name)
		VALUES 
			('Алхимия'),
			('Ангелы'),
			('Демоны'),
			('Драконы'),
			('Зомби'),
			('Вампиры'),
			('Геймеры'),
			('Герои'),
			('Гильдии'),
			('ГГ имба'),
			('ГГ женщина'),
			('ГГ мужчина'),
			('Выживание'),
			('Дружба'),
			('Империя'),
			('Королевство'),
			('Магия')
		ON CONFLICT (name) DO NOTHING;
	`

	if _, err := db.ExecContext(ctx, categoriesQuery); err != nil {
		log.Printf("Warning: failed to seed categories: %v", err)
	}
	
	return nil
}