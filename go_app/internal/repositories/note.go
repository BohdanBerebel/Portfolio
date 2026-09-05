package repositories

import (
	"context"
	"errors"

	"github.com/BohdanBerebel/Portfolio/go_app/internal/models"

	apperrors "github.com/BohdanBerebel/Portfolio/go_app/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepository struct {
	db *pgxpool.Pool
}

func NewNoteRepository(db *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{
		db: db,
	}
}

func (r *NoteRepository) Create(
	ctx context.Context,
	note *models.Note,
) error {
	query := `
		INSERT INTO notes (
			user_id,
			title,
			content
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		note.UserID,
		note.Title,
		note.Content,
	).Scan(
		&note.ID,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
}

func (r *NoteRepository) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Note, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]models.Note, 0)

	for rows.Next() {
		var note models.Note

		if err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		); err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NoteRepository) GetByID(
	ctx context.Context,
	noteID int64,
	userID int64,
) (*models.Note, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE id = $1 AND user_id = $2
	`

	var note models.Note

	err := r.db.QueryRow(
		ctx,
		query,
		noteID,
		userID,
	).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNoteNotFound
		}

		return nil, err
	}

	return &note, nil
}

func (r *NoteRepository) Update(
	ctx context.Context,
	noteID int64,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	query := `
		UPDATE notes
		SET title = $1,
		    content = $2,
		    updated_at = NOW()
		WHERE id = $3
		  AND user_id = $4
		RETURNING id, user_id, title, content, created_at, updated_at
	`

	var note models.Note

	err := r.db.QueryRow(
		ctx,
		query,
		title,
		content,
		noteID,
		userID,
	).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNoteNotFound
		}

		return nil, err
	}

	return &note, nil
}

func (r *NoteRepository) Delete(
	ctx context.Context,
	noteID int64,
	userID int64,
) error {
	query := `
		DELETE FROM notes
		WHERE id = $1
		  AND user_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		noteID,
		userID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrNoteNotFound
	}

	return nil
}
