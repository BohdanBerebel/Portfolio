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

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	databaseURL := os.Getenv("TEST_DATABASE_URL")

	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(
		ctx,
		databaseURL,
	)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestNoteRepository_Create(t *testing.T) {
	db := setupTestDB(t)

	repository := NewNoteRepository(db)

	ctx := context.Background()

	var userID int64

	err := db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id`,
		"integration-test@example.com",
		"dummy-hash",
	).Scan(&userID)

	if err != nil {
		t.Fatalf(
			"failed to create test user: %v",
			err,
		)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM notes WHERE user_id = $1",
			userID,
		)

		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM users WHERE id = $1",
			userID,
		)
	})

	note := &models.Note{
		UserID:  userID,
		Title:   "Integration test",
		Content: "Hello PostgreSQL",
	}

	err = repository.Create(
		ctx,
		note,
	)
	if err != nil {
		t.Fatalf(
			"Create failed: %v",
			err,
		)
	}

	if note.ID == 0 {
		t.Fatal("expected note ID to be populated")
	}

	if note.Title != "Integration test" {
		t.Fatalf("unexpected title: %q", note.Title)
	}

	if note.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be populated")
	}

	if note.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be populated")
	}

	var count int

	err = db.QueryRow(
		ctx,
		`SELECT COUNT(*)
		FROM notes
		WHERE id = $1`,
		note.ID,
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to verify note: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected note to exist, got count %d", count)
	}

}

func TestNoteRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repository := NewNoteRepository(db)

	ctx := context.Background()

	var userID int64

	err := db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id`,
		"get-by-id-test@example.com",
		"dummy-hash",
	).Scan(&userID)

	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM notes WHERE user_id = $1",
			userID,
		)

		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM users WHERE id = $1",
			userID,
		)
	})

	note := &models.Note{
		UserID:  userID,
		Title:   "GetByID test",
		Content: "Testing repository",
	}

	err = repository.Create(ctx, note)
	if err != nil {
		t.Fatalf("failed to create test note: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		result, err := repository.GetByID(
			ctx,
			note.ID,
			userID,
		)

		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if result.ID != note.ID {
			t.Fatalf(
				"expected ID %d, got %d",
				note.ID,
				result.ID,
			)
		}

		if result.UserID != userID {
			t.Fatalf(
				"expected user ID %d, got %d",
				userID,
				result.UserID,
			)
		}

		if result.Title != note.Title {
			t.Fatalf(
				"expected title %q, got %q",
				note.Title,
				result.Title,
			)
		}

		if result.Content != note.Content {
			t.Fatalf(
				"expected content %q, got %q",
				note.Content,
				result.Content,
			)
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		_, err := repository.GetByID(
			ctx,
			note.ID,
			userID+1,
		)

		if !errors.Is(err, apperrors.ErrNoteNotFound) {
			t.Fatalf(
				"expected ErrNoteNotFound, got %v",
				err,
			)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repository.GetByID(
			ctx,
			999999999,
			userID,
		)

		if !errors.Is(err, apperrors.ErrNoteNotFound) {
			t.Fatalf(
				"expected ErrNoteNotFound, got %v",
				err,
			)
		}
	})
}

func TestNoteRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repository := NewNoteRepository(db)

	ctx := context.Background()

	var userA int64
	var userB int64

	err := db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id`,
		"get-all-user-a@example.com",
		"dummy-hash",
	).Scan(&userA)

	if err != nil {
		t.Fatalf("failed to create user A: %v", err)
	}

	err = db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id`,
		"get-all-user-b@example.com",
		"dummy-hash",
	).Scan(&userB)

	if err != nil {
		t.Fatalf("failed to create user B: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM notes WHERE user_id IN ($1, $2)",
			userA,
			userB,
		)

		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM users WHERE id IN ($1, $2)",
			userA,
			userB,
		)
	})

	notesA := []*models.Note{
		{
			UserID:  userA,
			Title:   "A note 1",
			Content: "Content A1",
		},
		{
			UserID:  userA,
			Title:   "A note 2",
			Content: "Content A2",
		},
	}

	noteB := &models.Note{
		UserID:  userB,
		Title:   "B note 1",
		Content: "Content B1",
	}

	for _, note := range notesA {
		if err := repository.Create(ctx, note); err != nil {
			t.Fatalf("failed to create user A note: %v", err)
		}
	}

	if err := repository.Create(ctx, noteB); err != nil {
		t.Fatalf("failed to create user B note: %v", err)
	}

	result, err := repository.GetByUserID(ctx, userA)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 notes, got %d",
			len(result),
		)
	}

	for _, note := range result {
		if note.UserID != userA {
			t.Fatalf(
				"expected note to belong to user %d, got %d",
				userA,
				note.UserID,
			)
		}
	}

	if result[0].Title != "A note 2" {
		t.Fatalf(
			"expected first note to be %q, got %q",
			"A note 2",
			result[0].Title,
		)
	}

	if result[1].Title != "A note 1" {
		t.Fatalf(
			"expected second note to be %q, got %q",
			"A note 1",
			result[1].Title,
		)
	}
}

func TestNoteRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repository := NewNoteRepository(db)

	ctx := context.Background()

	var userID int64

	err := db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id`,
		"update-test@example.com",
		"dummy-hash",
	).Scan(&userID)

	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM notes WHERE user_id = $1",
			userID,
		)

		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM users WHERE id = $1",
			userID,
		)
	})

	note := &models.Note{
		UserID:  userID,
		Title:   "Original title",
		Content: "Original content",
	}

	err = repository.Create(ctx, note)
	if err != nil {
		t.Fatalf("failed to create test note: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		updated, err := repository.Update(
			ctx,
			note.ID,
			userID,
			"Updated title",
			"Updated content",
		)

		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if updated.ID != note.ID {
			t.Fatalf(
				"expected ID %d, got %d",
				note.ID,
				updated.ID,
			)
		}

		if updated.UserID != userID {
			t.Fatalf(
				"expected user ID %d, got %d",
				userID,
				updated.UserID,
			)
		}

		if updated.Title != "Updated title" {
			t.Fatalf(
				"expected title %q, got %q",
				"Updated title",
				updated.Title,
			)
		}

		if updated.Content != "Updated content" {
			t.Fatalf(
				"expected content %q, got %q",
				"Updated content",
				updated.Content,
			)
		}

		if updated.UpdatedAt.IsZero() {
			t.Fatal("expected UpdatedAt to be populated")
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		_, err := repository.Update(
			ctx,
			note.ID,
			userID+1,
			"Hacked title",
			"Hacked content",
		)

		if !errors.Is(err, apperrors.ErrNoteNotFound) {
			t.Fatalf(
				"expected ErrNoteNotFound, got %v",
				err,
			)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repository.Update(
			ctx,
			999999999,
			userID,
			"Updated",
			"Updated",
		)

		if !errors.Is(err, apperrors.ErrNoteNotFound) {
			t.Fatalf(
				"expected ErrNoteNotFound, got %v",
				err,
			)
		}
	})
}

func TestNoteRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repository := NewNoteRepository(db)

	ctx := context.Background()

	var userID int64

	err := db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id`,
		"delete-test@example.com",
		"dummy-hash",
	).Scan(&userID)

	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM notes WHERE user_id = $1",
			userID,
		)

		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM users WHERE id = $1",
			userID,
		)
	})

	note := &models.Note{
		UserID:  userID,
		Title:   "Delete test",
		Content: "This will be deleted",
	}

	err = repository.Create(ctx, note)
	if err != nil {
		t.Fatalf("failed to create test note: %v", err)
	}

	t.Run("wrong user", func(t *testing.T) {
		err := repository.Delete(
			ctx,
			note.ID,
			userID+1,
		)

		if !errors.Is(err, apperrors.ErrNoteNotFound) {
			t.Fatalf(
				"expected ErrNoteNotFound, got %v",
				err,
			)
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repository.Delete(
			ctx,
			999999999,
			userID,
		)

		if !errors.Is(err, apperrors.ErrNoteNotFound) {
			t.Fatalf(
				"expected ErrNoteNotFound, got %v",
				err,
			)
		}
	})

	t.Run("success", func(t *testing.T) {
		err := repository.Delete(
			ctx,
			note.ID,
			userID,
		)

		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = repository.GetByID(
			ctx,
			note.ID,
			userID,
		)

		if !errors.Is(err, apperrors.ErrNoteNotFound) {
			t.Fatalf(
				"expected deleted note to be not found, got %v",
				err,
			)
		}
	})
}
