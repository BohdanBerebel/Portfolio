package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/BohdanBerebel/Portfolio/go_app/internal/errors"
	"github.com/BohdanBerebel/Portfolio/go_app/pkg/response"
)

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		response.Error(
			c,
			http.StatusConflict,
			"user already exists",
		)

	case errors.Is(err, apperrors.ErrNoteNotFound):
		response.Error(
			c,
			http.StatusNotFound,
			"note not found",
		)

	case errors.Is(err, apperrors.ErrUserNotFound):
		response.Error(
			c,
			http.StatusNotFound,
			"user not found",
		)

	case errors.Is(err, apperrors.ErrInvalidCredentials):
		response.Error(
			c,
			http.StatusUnauthorized,
			"invalid credentials",
		)

	default:
		response.Error(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}
