package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type requestJSONChirp struct {
	RequestedChirp string `json:"body"`
}

type successfulResponseParams struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
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
		requestedChirpBadWordsRemoved := badWordsRemoved(requestJson)

		respBody := successfulResponseParams{
			CleanedBody: requestedChirpBadWordsRemoved,
		}
		respondWithJSON(w, 200, respBody)
	}
}
