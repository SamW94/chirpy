package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SamW94/chirpy/internal/auth"
	"github.com/google/uuid"
)

type requestJSONLogin struct {
	Password string        `json:"password"`
	Email    string        `json:"email"`
	Expiry   time.Duration `json:"expires_in_seconds"`
}

type LoginSuccessfulResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	JWT       string    `json:"token"`
}

func (acfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := requestJSONLogin{}
	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Printf("Error decoding JSON from request body: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	if requestJson.Expiry < time.Second*3600 || requestJson.Expiry == 0 {
		requestJson.Expiry = time.Second * 3600
	}

	user, err := acfg.DatabaseQueries.RetrieveUserByEmail(context.Background(), requestJson.Email)
	if err != nil {
		log.Printf("Failed to retrieve user from database: %v", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, requestJson.Password)
	if err != nil {
		log.Printf("Call to CheckPasswordHash returned an error: %v", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	JWToken, err := auth.MakeJWT(user.ID, acfg.JWTSecret, time.Second*requestJson.Expiry)
	if err != nil {
		log.Printf("Error creating JWT for user ID %s: %v", user.ID, err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	respBody := LoginSuccessfulResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		JWT:       JWToken,
	}

	respondWithJSON(w, 200, respBody)
}
