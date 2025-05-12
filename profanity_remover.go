package main

import (
	"log"
	"slices"
	"strings"
)

func badWordsRemoved(r requestJSONChirp) string {
	log.Printf("Checking chirp body for profanity...")
	wordsInChirp := strings.Split(r.RequestedChirp, " ")
	cleanedWords := make([]string, 0)
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range wordsInChirp {
		if slices.Contains(badWords, strings.ToLower(word)) {
			log.Printf("Profanity found! %s :o", word)
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	return strings.Join(cleanedWords, " ")
}
