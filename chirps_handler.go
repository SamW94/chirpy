package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SamW94/chirpy/internal/database"
	"github.com/google/uuid"
)

type requestJSONChirp struct {
	RequestedChirp string    `json:"body"`
	UserID         uuid.UUID `json:"user_id"`
}

type successfulResponseParams struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (acfg *apiConfig) chirpsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := requestJSONChirp{}
	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Printf("Error decoding request body to requestJSON: %v", err)
		respondWithError(w, 500, "Something went wrong")
	}

	if len(requestJson.RequestedChirp) > 140 {
		log.Printf("The Chirp in the request body is too long (more than 140 character), failing with HTTP 400...")
		respondWithError(w, 400, "Chirp is too long")
	} else {
		createChirpParams := database.CreateChirpParams{
			Body:   badWordsRemoved(requestJson),
			UserID: requestJson.UserID,
		}

		chirp, err := acfg.DatabaseQueries.CreateChirp(context.Background(), createChirpParams)
		if err != nil {
			log.Printf("Error writing chirp to database: %v", err)
		}

		log.Printf("Successfully created chirp with ID %v for user with ID %v", chirp.ID, chirp.UserID)
		respBody := successfulResponseParams{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		respondWithJSON(w, 201, respBody)
	}
}
