package services

import (
	"context"

	"github.com/BohdanBerebel/go_app/internal/auth"
	apperrors "github.com/BohdanBerebel/go_app/internal/errors"
	"github.com/BohdanBerebel/go_app/internal/models"
)

type UserService struct {
	repository UserRepository
	jwtSecret  string
}

type AuthResult struct {
	User  *models.User
	Token string
}

type UserRepository interface {
	Create(
		ctx context.Context,
		user *models.User,
	) error

	FindByEmail(
		ctx context.Context,
		email string,
	) (*models.User, error)
}

func NewUserService(
	repository UserRepository,
	jwtSecret string,
) *UserService {
	return &UserService{
		repository: repository,
		jwtSecret:  jwtSecret,
	}
}

func (s *UserService) Register(
	ctx context.Context,
	email string,
	password string,
) (*AuthResult, error) {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(
		user.ID,
		s.jwtSecret,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) Login(
	ctx context.Context,
	email string,
	password string,
) (*AuthResult, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	if err := auth.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(
		user.ID,
		s.jwtSecret,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:  user,
		Token: token,
	}, nil
}
