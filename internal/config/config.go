package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	domain_evently "evently/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func MigrateDB(cfg domain_evently.DBConfig) error {
	var dsn string
	if cfg.URL != "" {
		dsn = cfg.URL
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
		)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, cfg.MigrationsDir); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}

func LoadConfig() *domain_evently.Config {
	return &domain_evently.Config{
		DB: domain_evently.DBConfig{
			URL:           getEnv("DATABASE_URL", ""),
			Host:          getEnv("DB_HOST", "localhost"),
			Port:          getEnv("DB_PORT", "5432"),
			User:          getEnv("DB_USER", "arman"),
			Password:      getEnv("DB_PASSWORD", ""),
			DBName:        getEnv("DB_NAME", "postgres"),
			SSLMode:       getEnv("DB_SSLMODE", "disable"),
			MigrationsDir: getEnv("MIGRATIONS_DIR", "../migrations"),
		},
		JWT: domain_evently.JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
		},
	}
}

func NewPGXPool(ctx context.Context, cfg domain_evently.DBConfig) (*pgxpool.Pool, error) {
	var dsn string
	if cfg.URL != "" {
		dsn = cfg.URL
	} else {
		if cfg.Host == "" || cfg.User == "" || cfg.DBName == "" {
			return nil, fmt.Errorf("incomplete DB config: supply DATABASE_URL or PGHOST, PGUSER, PGPASSWORD, PGDATABASE")
		}
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
		)
	}

	pcfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	if err := MigrateDB(cfg); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	return pool, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
