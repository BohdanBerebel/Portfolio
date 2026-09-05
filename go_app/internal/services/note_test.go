package services

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/BohdanBerebel/Portfolio/go_app/internal/errors"
	"github.com/BohdanBerebel/Portfolio/go_app/internal/models"
)

type fakeNoteRepository struct {
	createdNote *models.Note
	createErr   error

	note   *models.Note
	getErr error
}

func (f *fakeNoteRepository) Create(
	ctx context.Context,
	note *models.Note,
) error {
	f.createdNote = note

	return f.createErr
}

func (f *fakeNoteRepository) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Note, error) {
	return nil, nil
}

func (f *fakeNoteRepository) GetByID(
	ctx context.Context,
	noteID int64,
	userID int64,
) (*models.Note, error) {
	return f.note, f.getErr
}

func (f *fakeNoteRepository) Update(
	ctx context.Context,
	noteID int64,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	return nil, nil
}

func (f *fakeNoteRepository) Delete(
	ctx context.Context,
	noteID int64,
	userID int64,
) error {
	return nil
}

func TestNoteService_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		noteID         int64
		userID         int64
		repositoryNote *models.Note
		repositoryErr  error
		wantErr        error
	}{
		{
			name:   "success",
			noteID: 1,
			userID: 42,
			repositoryNote: &models.Note{
				ID:      1,
				UserID:  42,
				Title:   "My note",
				Content: "Hello",
			},
		},
		{
			name:          "note not found",
			noteID:        999,
			userID:        42,
			repositoryErr: apperrors.ErrNoteNotFound,
			wantErr:       apperrors.ErrNoteNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeNoteRepository{
				note:   tt.repositoryNote,
				getErr: tt.repositoryErr,
			}

			service := NewNoteService(repository)

			note, err := service.GetByID(
				context.Background(),
				tt.noteID,
				tt.userID,
			)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf(
						"expected error %v, got %v",
						tt.wantErr,
						err,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"unexpected error: %v",
					err,
				)
			}

			if note == nil {
				t.Fatal("expected note, got nil")
			}

			if note.ID != tt.repositoryNote.ID {
				t.Fatalf(
					"expected note ID %d, got %d",
					tt.repositoryNote.ID,
					note.ID,
				)
			}
		})
	}
}
