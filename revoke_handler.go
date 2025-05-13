package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

func (acfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
	_, err := acfg.DatabaseQueries.UpdateRefreshTokenRevoke(context.Background(), refreshTokenString)
	if err != nil {
		log.Printf("Error updating refresh token in database - refresh token may not have been revoked: %v", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
