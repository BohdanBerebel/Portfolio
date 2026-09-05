package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/BohdanBerebel/go_app/internal/errors"
	"github.com/BohdanBerebel/go_app/internal/models"
	"github.com/gin-gonic/gin"
)

type fakeNoteService struct {
	note  *models.Note
	notes []models.Note

	createErr error
	getErr    error
	updateErr error
	deleteErr error
}

func (f *fakeNoteService) Create(
	ctx context.Context,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}

	return f.note, nil
}

func (f *fakeNoteService) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Note, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}

	return f.notes, nil
}

func (f *fakeNoteService) GetByID(
	ctx context.Context,
	noteID int64,
	userID int64,
) (*models.Note, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}

	return f.note, nil
}

func (f *fakeNoteService) Update(
	ctx context.Context,
	noteID int64,
	userID int64,
	title string,
	content string,
) (*models.Note, error) {
	if f.updateErr != nil {
		return nil, f.updateErr
	}

	return f.note, nil
}

func (f *fakeNoteService) Delete(
	ctx context.Context,
	noteID int64,
	userID int64,
) error {
	return f.deleteErr
}

func setTestUserID(userID int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func TestNoteHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeNoteService{
		notes: []models.Note{
			{
				ID:      1,
				UserID:  42,
				Title:   "My note",
				Content: "Hello",
			},
		},
	}

	handler := NewNoteHandler(service)

	router := gin.New()

	router.GET("/notes", setTestUserID(42), handler.GetAll)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var response struct {
		Data []models.Note `json:"data"`
	}

	err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	)
	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if len(response.Data) != 1 {
		t.Fatalf(
			"expected 1 note, got %d",
			len(response.Data),
		)
	}

	if response.Data[0].ID != 1 {
		t.Fatalf(
			"expected note ID 1, got %d",
			response.Data[0].ID,
		)
	}
}

func TestNoteHandler_GetAll_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeNoteService{
		getErr: errors.New("database error"),
	}

	handler := NewNoteHandler(service)

	router := gin.New()

	router.GET("/notes", setTestUserID(42), handler.GetAll)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}
}

func TestNoteHandler_Create(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		authenticated bool
		serviceErr    error
		wantStatus    int
	}{
		{
			name:          "success",
			body:          `{"title":"My note","content":"Hello"}`,
			authenticated: true,
			wantStatus:    http.StatusCreated,
		},
		{
			name:          "invalid json",
			body:          `{"title":`,
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name:          "validation error",
			body:          `{"title":"","content":"Hello"}`,
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name:          "service error",
			body:          `{"title":"My note","content":"Hello"}`,
			authenticated: true,
			serviceErr:    errors.New("database error"),
			wantStatus:    http.StatusInternalServerError,
		},
		{
			name:          "unauthorized",
			body:          `{"title":"My note","content":"Hello"}`,
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeNoteService{
				note: &models.Note{
					ID:      1,
					UserID:  42,
					Title:   "My note",
					Content: "Hello",
				},
				createErr: tt.serviceErr,
			}

			handler := NewNoteHandler(service)

			router := gin.New()

			if tt.authenticated {
				router.POST(
					"/notes",
					setTestUserID(42),
					handler.Create,
				)
			} else {
				router.POST(
					"/notes",
					handler.Create,
				)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/notes",
				bytes.NewBufferString(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					recorder.Code,
				)
			}
			if tt.wantStatus == http.StatusCreated {
				var response struct {
					Data models.Note `json:"data"`
				}

				if err := json.Unmarshal(
					recorder.Body.Bytes(),
					&response,
				); err != nil {
					t.Fatalf(
						"failed to decode response: %v",
						err,
					)
				}

				if response.Data.ID != 1 {
					t.Fatalf(
						"expected ID 1, got %d",
						response.Data.ID,
					)
				}
			}
		})
	}
}

func TestNoteHandler_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		noteID        string
		authenticated bool
		serviceErr    error
		wantStatus    int
	}{
		{
			name:          "success",
			noteID:        "1",
			authenticated: true,
			wantStatus:    http.StatusOK,
		},
		{
			name:          "invalid id",
			noteID:        "abc",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name:          "not found",
			noteID:        "999",
			authenticated: true,
			serviceErr:    apperrors.ErrNoteNotFound,
			wantStatus:    http.StatusNotFound,
		},
		{
			name:          "internal error",
			noteID:        "1",
			authenticated: true,
			serviceErr:    errors.New("database error"),
			wantStatus:    http.StatusInternalServerError,
		},
		{
			name:          "unauthorized",
			noteID:        "1",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeNoteService{
				note: &models.Note{
					ID:      1,
					UserID:  42,
					Title:   "My note",
					Content: "Hello",
				},
				getErr: tt.serviceErr,
			}

			handler := NewNoteHandler(service)

			router := gin.New()

			if tt.authenticated {
				router.GET(
					"/notes/:id",
					setTestUserID(42),
					handler.GetByID,
				)
			} else {
				router.GET(
					"/notes/:id",
					handler.GetByID,
				)
			}

			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/notes/%s", tt.noteID),
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					recorder.Code,
				)
			}
			if tt.wantStatus == http.StatusOK {
				var response struct {
					Data models.Note `json:"data"`
				}

				if err := json.Unmarshal(
					recorder.Body.Bytes(),
					&response,
				); err != nil {
					t.Fatalf(
						"failed to decode response: %v",
						err,
					)
				}

				if response.Data.ID != 1 {
					t.Fatalf(
						"expected ID 1, got %d",
						response.Data.ID,
					)
				}
			}
		})
	}
}

func TestNoteHandler_Update(t *testing.T) {
	tests := []struct {
		name          string
		noteID        string
		body          string
		authenticated bool
		serviceErr    error
		wantStatus    int
	}{
		{
			name:          "success",
			noteID:        "1",
			body:          `{"title":"Updated","content":"New content"}`,
			authenticated: true,
			wantStatus:    http.StatusOK,
		},
		{
			name:          "invalid id",
			noteID:        "abc",
			body:          `{"title":"Updated","content":"New content"}`,
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name:          "invalid json",
			noteID:        "1",
			body:          `{"title":`,
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name:          "not found",
			noteID:        "999",
			body:          `{"title":"Updated","content":"New content"}`,
			authenticated: true,
			serviceErr:    apperrors.ErrNoteNotFound,
			wantStatus:    http.StatusNotFound,
		},
		{
			name:          "internal error",
			noteID:        "1",
			body:          `{"title":"Updated","content":"New content"}`,
			authenticated: true,
			serviceErr:    errors.New("database error"),
			wantStatus:    http.StatusInternalServerError,
		},
		{
			name:          "unauthorized",
			noteID:        "1",
			body:          `{"title":"Updated","content":"New content"}`,
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeNoteService{
				note: &models.Note{
					ID:      1,
					UserID:  42,
					Title:   "Updated",
					Content: "New content",
				},
				updateErr: tt.serviceErr,
			}

			handler := NewNoteHandler(service)

			router := gin.New()

			if tt.authenticated {
				router.PUT(
					"/notes/:id",
					setTestUserID(42),
					handler.Update,
				)
			} else {
				router.PUT(
					"/notes/:id",
					handler.Update,
				)
			}

			req := httptest.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/notes/%s", tt.noteID),
				bytes.NewBufferString(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					recorder.Code,
				)
			}

			if tt.wantStatus == http.StatusOK {
				var response struct {
					Data models.Note `json:"data"`
				}

				if err := json.Unmarshal(
					recorder.Body.Bytes(),
					&response,
				); err != nil {
					t.Fatalf(
						"failed to decode response: %v",
						err,
					)
				}

				if response.Data.ID != 1 {
					t.Fatalf(
						"expected ID 1, got %d",
						response.Data.ID,
					)
				}
			}
		})
	}
}

func TestNoteHandler_Delete(t *testing.T) {
	tests := []struct {
		name          string
		noteID        string
		authenticated bool
		serviceErr    error
		wantStatus    int
	}{
		{
			name:          "success",
			noteID:        "1",
			authenticated: true,
			wantStatus:    http.StatusNoContent,
		},
		{
			name:          "invalid id",
			noteID:        "abc",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name:          "not found",
			noteID:        "999",
			authenticated: true,
			serviceErr:    apperrors.ErrNoteNotFound,
			wantStatus:    http.StatusNotFound,
		},
		{
			name:          "internal error",
			noteID:        "1",
			authenticated: true,
			serviceErr:    errors.New("database error"),
			wantStatus:    http.StatusInternalServerError,
		},
		{
			name:          "unauthorized",
			noteID:        "1",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeNoteService{
				deleteErr: tt.serviceErr,
			}

			handler := NewNoteHandler(service)

			router := gin.New()

			if tt.authenticated {
				router.DELETE(
					"/notes/:id",
					setTestUserID(42),
					handler.Delete,
				)
			} else {
				router.DELETE(
					"/notes/:id",
					handler.Delete,
				)
			}

			req := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/notes/%s", tt.noteID),
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					recorder.Code,
				)
			}

			if tt.wantStatus == http.StatusNoContent {
				if recorder.Body.Len() != 0 {
					t.Fatalf(
						"expected empty response body, got %q",
						recorder.Body.String(),
					)
				}
			}
		})
	}
}
