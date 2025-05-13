package auth

import (
	"testing"
)

func TestMakeRefreshToken_Success(t *testing.T) {
	refreshToken, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("MakeRefreshToken failed: %v", err)
	}

	if len(refreshToken) < 32 {
		t.Errorf("expected 32 characters , got %d", len(refreshToken))
	}
}
