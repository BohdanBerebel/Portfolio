package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestHashPassword(t *testing.T) {
	password := "super-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == password {
		t.Fatal("password must not be stored as plaintext")
	}

	if err := CheckPassword(password, hash); err != nil {
		t.Fatalf("expected password to match: %v", err)
	}
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = CheckPassword("wrong-password", hash)

	if err == nil {
		t.Fatal("expected password check to fail")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	const secret = "test-secret"

	token, err := GenerateToken(42, secret)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	userID, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if userID != 42 {
		t.Fatalf(
			"expected user ID 42, got %d",
			userID,
		)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken(
		42,
		"correct-secret",
	)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ValidateToken(
		token,
		"wrong-secret",
	)

	if err == nil {
		t.Fatal("expected validation to fail with wrong secret")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := ValidateToken(
		"this-is-not-a-jwt",
		"test-secret",
	)

	if err == nil {
		t.Fatal("expected invalid token error")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	const secret = "test-secret"

	claims := Claims{
		UserID: 42,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(-time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(
				time.Now().Add(-2 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(secret),
	)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = ValidateToken(
		tokenString,
		secret,
	)

	if err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}

func TestUserIDContext(t *testing.T) {
	c, _ := gin.CreateTestContext(
		httptest.NewRecorder(),
	)

	SetUserID(c, 42)

	userID, ok := GetUserID(c)

	if !ok {
		t.Fatal("expected user ID to exist")
	}

	if userID != 42 {
		t.Fatalf(
			"expected user ID 42, got %d",
			userID,
		)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	c, _ := gin.CreateTestContext(
		httptest.NewRecorder(),
	)

	userID, ok := GetUserID(c)

	if ok {
		t.Fatal("expected user ID to be missing")
	}

	if userID != 0 {
		t.Fatalf(
			"expected user ID 0, got %d",
			userID,
		)
	}
}
