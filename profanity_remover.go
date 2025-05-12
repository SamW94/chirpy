package main

import (
	"slices"
	"strings"
)

func badWordsRemoved(r requestJSONChirp) string {
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
