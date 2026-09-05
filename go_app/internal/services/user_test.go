package services

import (
	"context"
	"errors"
	"testing"

	"github.com/BohdanBerebel/Portfolio/go_app/internal/auth"
	apperrors "github.com/BohdanBerebel/Portfolio/go_app/internal/errors"
	"github.com/BohdanBerebel/Portfolio/go_app/internal/models"
)

type fakeUserRepository struct {
	user      *models.User
	createErr error
	findErr   error
}

func (f *fakeUserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {
	f.user = user

	return f.createErr
}

func (f *fakeUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}

	return f.user, nil
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		repositoryErr error
		wantErr       error
	}{
		{
			name:     "success",
			email:    "test@example.com",
			password: "super-secret-password",
		},
		{
			name:          "repository error",
			email:         "test@example.com",
			password:      "super-secret-password",
			repositoryErr: errors.New("database error"),
			wantErr:       errors.New("database error"),
		},
		{
			name:          "user already exists",
			email:         "test@example.com",
			password:      "super-secret-password",
			repositoryErr: apperrors.ErrUserAlreadyExists,
			wantErr:       apperrors.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeUserRepository{
				createErr: tt.repositoryErr,
			}

			service := NewUserService(
				repository,
				"test-secret",
			)

			result, err := service.Register(
				context.Background(),
				tt.email,
				tt.password,
			)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf(
						"expected error %v, got nil",
						tt.wantErr,
					)
				}

				if tt.wantErr.Error() != err.Error() {
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

			if result == nil {
				t.Fatal("expected auth result, got nil")
			}

			if result.User == nil {
				t.Fatal("expected user, got nil")
			}

			if result.User.Email != tt.email {
				t.Fatalf(
					"expected email %q, got %q",
					tt.email,
					result.User.Email,
				)
			}

			if result.User.PasswordHash == "" {
				t.Fatal("expected password hash to be populated")
			}

			if result.User.PasswordHash == tt.password {
				t.Fatal("password must not be stored as plaintext")
			}

			if result.Token == "" {
				t.Fatal("expected token to be generated")
			}

			if repository.user == nil {
				t.Fatal("expected repository to receive user")
			}

			if repository.user.Email != tt.email {
				t.Fatalf(
					"expected repository email %q, got %q",
					tt.email,
					repository.user.Email,
				)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	passwordHash, err := auth.HashPassword("correct-password")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	tests := []struct {
		name          string
		email         string
		password      string
		repositoryErr error
		user          *models.User
		wantErr       error
	}{
		{
			name:     "success",
			email:    "test@example.com",
			password: "correct-password",
			user: &models.User{
				ID:           42,
				Email:        "test@example.com",
				PasswordHash: passwordHash,
			},
		},
		{
			name:          "user not found",
			email:         "test@example.com",
			password:      "correct-password",
			repositoryErr: apperrors.ErrUserNotFound,
			wantErr:       apperrors.ErrInvalidCredentials,
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrong-password",
			user: &models.User{
				ID:           42,
				Email:        "test@example.com",
				PasswordHash: passwordHash,
			},
			wantErr: apperrors.ErrInvalidCredentials,
		},
		{
			name:          "repository error",
			email:         "test@example.com",
			password:      "correct-password",
			repositoryErr: errors.New("database error"),
			wantErr:       apperrors.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeUserRepository{
				user:    tt.user,
				findErr: tt.repositoryErr,
			}

			service := NewUserService(
				repository,
				"test-secret",
			)

			result, err := service.Login(
				context.Background(),
				tt.email,
				tt.password,
			)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf(
						"expected error %v, got nil",
						tt.wantErr,
					)
				}

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

			if result == nil {
				t.Fatal("expected login result, got nil")
			}

			if result.User == nil {
				t.Fatal("expected user, got nil")
			}

			if result.User.ID != 42 {
				t.Fatalf(
					"expected user ID 42, got %d",
					result.User.ID,
				)
			}

			if result.User.Email != tt.email {
				t.Fatalf(
					"expected email %q, got %q",
					tt.email,
					result.User.Email,
				)
			}

			if result.Token == "" {
				t.Fatal("expected token to be generated")
			}
		})
	}
}
