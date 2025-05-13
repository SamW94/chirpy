package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT_Success(t *testing.T) {
	secret := "test_secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("expected userID %v, got %v", userID, validatedUserID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	secret := "test_secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil || !strings.Contains(err.Error(), "token is expired") {
		t.Errorf("expected expiration error, got %v", err)
	}
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	secret := "test_secret"
	wrongSecret := "wrong_secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil || !strings.Contains(err.Error(), "signature is invalid") {
		t.Errorf("expected signature error, got %v", err)
	}
}

func TestValidateJWT_InvalidUUID(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   "not-a-valid-uuid",
	})

	secret := "test_secret"
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("SignedString failed: %v", err)
	}

	_, err = ValidateJWT(signedToken, secret)
	if err == nil || !strings.Contains(err.Error(), "error parsing userID string") {
		t.Errorf("expected UUID parsing error, got %v", err)
	}
}
