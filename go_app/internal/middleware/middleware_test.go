package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/BohdanBerebel/Portfolio/go_app/internal/auth"
)

func setupAuthRouter(jwtSecret string) *gin.Engine {
	router := gin.New()

	router.Use(AuthMiddleware(jwtSecret))

	router.GET("/protected", func(c *gin.Context) {
		userID, ok := auth.GetUserID(c)

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "user ID not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
		})
	})

	return router
}

func TestAuthMiddleware_Success(t *testing.T) {
	const secret = "test-secret"

	token, err := auth.GenerateToken(42, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	router := setupAuthRouter(secret)

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
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
		UserID int64 `json:"user_id"`
	}

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.UserID != 42 {
		t.Fatalf(
			"expected user ID 42, got %d",
			response.UserID,
		)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	router := setupAuthRouter("test-secret")

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			recorder.Code,
		)
	}
}

func TestAuthMiddleware_InvalidHeader(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "missing bearer",
			header: "some-token",
		},
		{
			name:   "wrong scheme",
			header: "Basic some-token",
		},
		{
			name:   "too many parts",
			header: "Bearer token extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupAuthRouter("test-secret")

			req := httptest.NewRequest(
				http.MethodGet,
				"/protected",
				nil,
			)

			req.Header.Set(
				"Authorization",
				tt.header,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf(
					"expected status %d, got %d",
					http.StatusUnauthorized,
					recorder.Code,
				)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := setupAuthRouter("test-secret")

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer definitely-not-a-jwt",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			recorder.Code,
		)
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	token, err := auth.GenerateToken(
		42,
		"correct-secret",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	router := setupAuthRouter("wrong-secret")

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			recorder.Code,
		)
	}
}

func TestAuthMiddleware_DoesNotCallNextOnFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	handlerCalled := false

	router.Use(AuthMiddleware("test-secret"))

	router.GET("/protected", func(c *gin.Context) {
		handlerCalled = true

		c.JSON(http.StatusOK, gin.H{
			"status": "should not happen",
		})
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			recorder.Code,
		)
	}

	if handlerCalled {
		t.Fatal("protected handler should not be called")
	}
}
