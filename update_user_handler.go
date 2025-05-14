package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SamW94/chirpy/internal/auth"
	"github.com/SamW94/chirpy/internal/database"
	"github.com/google/uuid"
)

type requestJSONUpdateUser struct {
	Password      string `json:"password"`
	RequestedUser string `json:"email"`
}

type UpdateUserResponseSuccessful struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (acfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := requestJSONUpdateUser{}
	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Printf("Error decoding request body to requestJSON: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Bearer token from HTTP request headers: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	userID, err := auth.ValidateJWT(token, acfg.JWTSecret)
	if err != nil {
		log.Printf("Error while validating JWT: %v", err)
		respondWithError(w, 401, "Invalid JWT")
		return
	}

	if requestJson.Password == "" {
		log.Printf("No password provided in request body.")
		respondWithError(w, 500, "No password provided in request body.")
		return
	}

	hashedPassword, err := auth.HashPassword(requestJson.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	updateUserParams := database.UpdateEmailAndPasswordParams{
		Email:          requestJson.RequestedUser,
		HashedPassword: hashedPassword,
		ID:             userID,
	}

	dbUser, err := acfg.DatabaseQueries.UpdateEmailAndPassword(context.Background(), updateUserParams)
	if err != nil {
		log.Printf("Error writing chirp to database: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	log.Printf("Successfully updated password and email for user with ID %v", dbUser.ID)
	respBody := UpdateUserResponseSuccessful{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, 200, respBody)
}
