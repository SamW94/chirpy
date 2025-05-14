package main

import (
	"context"
	"log"
	"net/http"

	"github.com/SamW94/chirpy/internal/auth"
	"github.com/SamW94/chirpy/internal/database"
	"github.com/google/uuid"
)

func (acfg *apiConfig) deleteChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirp_id"))
	log.Printf("Deleting chirp with ID: %v", chirpID)
	if err != nil {
		log.Printf("Error converting string in URL to UUID type for DB lookup: %v", err)
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
		respondWithError(w, 401, "Something went wrong")
		return
	}

	dbChirpID, err := acfg.DatabaseQueries.RetrieveChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error retreiving chirp from database: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if jwtUUID != dbChirpID.UserID {
		log.Printf("Current user does not match creator of chirp, will not delete")
		respondWithError(w, 403, "Forbidden")
		return
	}

	deleteChirpParams := database.DeleteChirpsParams{
		ID:     chirpID,
		UserID: jwtUUID,
	}

	err = acfg.DatabaseQueries.DeleteChirps(context.Background(), deleteChirpParams)
	if err != nil {
		log.Printf("Error deleting chirp %v from the database: %v", chirpID, err)
		respondWithError(w, 404, "Chirp not found")
		return
	}

	respondWithJSON(w, 204, nil)
}
