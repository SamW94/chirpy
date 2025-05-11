package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type requestJSON struct {
	RequestedChirp string `json:"body"`
}

type failedResponseParams struct {
	Error string `json:"error"`
}

type successfulResponseParams struct {
	CleanedBody string `json:"cleaned_body"`
}

func badWordsRemoved(r requestJSON) string {
	wordsInChirp := strings.Split(r.RequestedChirp, " ")
	cleanedWords := make([]string, 0)
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range wordsInChirp {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}

	return strings.Join(cleanedWords, " ")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respBody := failedResponseParams{
		Error: message,
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON to response body: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON to response body: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := requestJSON{}
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
