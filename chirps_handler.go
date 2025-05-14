package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SamW94/chirpy/internal/auth"
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

func (acfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := requestJSONChirp{}
	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Printf("Error decoding request body to requestJSON: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Bearer token from HTTP request headers: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	jwtUUID, err := auth.ValidateJWT(token, acfg.JWTSecret)
	if err != nil {
		log.Printf("Error while validating JWT: %v", err)
		respondWithError(w, 401, "Invalid JWT")
		return
	}

	if len(requestJson.RequestedChirp) > 140 {
		log.Printf("The Chirp in the request body is too long (more than 140 character), failing with HTTP 400...")
		respondWithError(w, 400, "Chirp is too long")
	} else {
		createChirpParams := database.CreateChirpParams{
			Body:   badWordsRemoved(requestJson),
			UserID: jwtUUID,
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

func (acfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("author_id") != "" {
		log.Printf("author_id query passed in URL, retrieving chirps for user with ID %v", r.URL.Query().Get("author_id"))
		userID, err := uuid.Parse(r.URL.Query().Get("author_id"))
		if err != nil {
			log.Printf("Error parsing value of author_id as UUID from string: %v", err)
			respondWithError(w, 500, "Something went wrong")
			return
		}

		chirps, err := acfg.DatabaseQueries.RetrieveChirpsByUserID(context.Background(), userID)
		if err != nil {
			log.Printf("Error retrieving chirps from the database: %v", err)
			respondWithError(w, 500, "Something went wrong")
			return
		}
		respondWithChirps(w, chirps)
		return
	}

	chirps, err := acfg.DatabaseQueries.RetrieveAllChirps(context.Background())
	if err != nil {
		log.Printf("Error retrieving chirps from the database: %v", err)
		return
	}
	respondWithChirps(w, chirps)
}

func (acfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirp_id"))
	log.Printf("Retrieving chirp with ID: %v", chirpID)
	if err != nil {
		log.Printf("Error converting string in URL to UUID type for DB lookup: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	chirp, err := acfg.DatabaseQueries.RetrieveChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error retrieving chirp by ID from database: %v", err)
		respondWithError(w, 404, fmt.Sprintf("Chirp with ID %v not found", chirpID))
	} else {
		log.Printf("Successfully retrieved chirp with ID %v from database", chirp.ID)
		respBody := successfulResponseParams{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		respondWithJSON(w, 200, respBody)
	}
}

func respondWithChirps(w http.ResponseWriter, chirps []database.Chirp) {
	respBody := []successfulResponseParams{}
	for _, chirp := range chirps {
		arrayItem := successfulResponseParams{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		respBody = append(respBody, arrayItem)
	}

	log.Printf("Responding with HTTP 200 and array of all chirps (oldest first) after successful creation of response body...")
	respondWithJSON(w, 200, respBody)
}
