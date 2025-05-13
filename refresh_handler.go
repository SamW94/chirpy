package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SamW94/chirpy/internal/auth"
)

type RefreshSuccessfulResponse struct {
	JWT string `json:"token"`
}

func (acfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)

	refreshTokenStruct, err := acfg.DatabaseQueries.RetrieveRefreshToken(context.Background(), refreshTokenString)
	if err != nil {
		log.Printf("Error retrieving refresh token from database: %v", err)
		respondWithError(w, 401, "Invalid refresh token.")
		return
	}

	if time.Now().After(refreshTokenStruct.ExpiresAt) || refreshTokenStruct.RevokedAt.Valid {
		log.Printf("Refresh token expired at %v - token is no longer valid", refreshTokenStruct.ExpiresAt)
		respondWithError(w, 401, "Invalid refresh token.")
		return
	}

	log.Println("Valid refresh token found in HTTP request.")
	userID, err := acfg.DatabaseQueries.GetUserFromRefreshToken(context.Background(), refreshTokenStruct.Token)
	if err != nil {
		log.Printf("Error retrieving user from refresh token: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	JWToken, err := auth.MakeJWT(userID, acfg.JWTSecret)
	if err != nil {
		log.Printf("Error creating JWT for user %v:\n %v", userID, err)
		respondWithError(w, 500, "Something went wrong.")
	}

	log.Printf("Successfully created JWT for userID %v", userID)
	respBody := RefreshSuccessfulResponse{
		JWT: JWToken,
	}

	respondWithJSON(w, 200, respBody)
}
