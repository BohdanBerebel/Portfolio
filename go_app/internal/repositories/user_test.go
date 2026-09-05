package repositories

import (
	"context"
	"errors"
	"os"
	"testing"

	apperrors "github.com/BohdanBerebel/Portfolio/go_app/internal/errors"
	"github.com/BohdanBerebel/Portfolio/go_app/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupUserTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	databaseURL := os.Getenv("TEST_DATABASE_URL")

	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		t.Fatalf("failed to ping database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestUserRepository_Create(t *testing.T) {
	db := setupUserTestDB(t)

	repository := NewUserRepository(db)

	ctx := context.Background()

	user := &models.User{
		Email:        "user-repository-test@example.com",
		PasswordHash: "dummy-hash",
	}

	err := repository.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM users WHERE id = $1",
			user.ID,
		)
	})

	if user.ID == 0 {
		t.Fatal("expected user ID to be populated")
	}

	if user.Email != "user-repository-test@example.com" {
		t.Fatalf(
			"unexpected email: %q",
			user.Email,
		)
	}

	if user.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be populated")
	}

	var count int

	err = db.QueryRow(
		ctx,
		`SELECT COUNT(*)
		FROM users
		WHERE id = $1`,
		user.ID,
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to verify user: %v", err)
	}

	if count != 1 {
		t.Fatalf(
			"expected user to exist, got count %d",
			count,
		)
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupUserTestDB(t)

	repository := NewUserRepository(db)

	ctx := context.Background()

	email := "user-repository-duplicate@example.com"

	user := &models.User{
		Email:        email,
		PasswordHash: "dummy-hash",
	}

	err := repository.Create(ctx, user)
	if err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM users WHERE id = $1",
			user.ID,
		)
	})

	duplicate := &models.User{
		Email:        email,
		PasswordHash: "another-hash",
	}

	err = repository.Create(ctx, duplicate)

	if !errors.Is(err, apperrors.ErrUserAlreadyExists) {
		t.Fatalf(
			"expected ErrUserAlreadyExists, got %v",
			err,
		)
	}
}
