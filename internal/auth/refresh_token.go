package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func MakeRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		log.Printf("Error filling refreshToken with random bytes using the rand.Read() function:\n %v", err)
		return "", fmt.Errorf("error creating random data from the crypto/rand package:\n %w", err)
	}

	return hex.EncodeToString(refreshToken), nil
}
