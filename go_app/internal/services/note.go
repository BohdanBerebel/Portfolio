package services

import (
	"context"

	"github.com/BohdanBerebel/go_app/internal/models"
)

type NoteService struct {
	repository NoteRepository
}

type NoteRepository interface {
	Create(
		ctx context.Context,
		note *models.Note,
	) error

	GetByUserID(
		ctx context.Context,
		userID int64,
	) ([]models.Note, error)

	GetByID(
		ctx context.Context,
		noteID int64,
		userID int64,
	) (*models.Note, error)

	Update(
		ctx context.Context,
		noteID int64,
		userID int64,
		title string,
		content string,
	) (*models.Note, error)

	Delete(
		ctx context.Context,
		noteID int64,
		userID int64,
	) error
}

func NewNoteService(
	repository NoteRepository,
) *NoteService {
	return &NoteService{
		repository: repository,
	}
}

func (s *NoteService) Create(
	ctx context.Context,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	note := &models.Note{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	if err := s.repository.Create(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Note, error) {
	return s.repository.GetByUserID(ctx, userID)
}

func (s *NoteService) GetByID(
	ctx context.Context,
	noteID int64,
	userID int64,
) (*models.Note, error) {
	return s.repository.GetByID(
		ctx,
		noteID,
		userID,
	)
}

func (s *NoteService) Update(
	ctx context.Context,
	noteID int64,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	return s.repository.Update(
		ctx,
		noteID,
		userID,
		title,
		content,
	)
}

func (s *NoteService) Delete(
	ctx context.Context,
	noteID int64,
	userID int64,
) error {
	return s.repository.Delete(
		ctx,
		noteID,
		userID,
	)
}
