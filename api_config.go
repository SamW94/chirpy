package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/SamW94/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits  atomic.Int32
	DatabaseQueries *database.Queries
	Platform        string
	JWTSecret       string
	PolkaKey        string
}

func (acfg *apiConfig) middlewareFileServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incrementing counter for: %s", r.URL.Path)
		acfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
