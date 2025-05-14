package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/SamW94/chirpy/internal/auth"
	"github.com/google/uuid"
)

type PolkaWebookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (acfg *apiConfig) polkaWebHookHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := PolkaWebookRequest{}
	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Printf("Error decoding JSON from request body: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	if requestJson.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting API key from HTTP request headers: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if key != acfg.PolkaKey {
		log.Printf("API key in request does not match the Polka API key required")
		respondWithError(w, 401, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(requestJson.Data.UserID)
	if err != nil {
		log.Printf("Error parsing user ID string in HTTP request to UUID object: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	_, err = acfg.DatabaseQueries.SetIsChirpyRedTrueForUser(context.Background(), userID)
	if err != nil {
		log.Printf("Error updating user in database: %v", err)
		respondWithError(w, 404, "Not Found")
		return
	}

	respondWithJSON(w, 204, nil)
}
