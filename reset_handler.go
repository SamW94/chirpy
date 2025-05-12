package main

import (
	"context"
	"log"
	"net/http"
)

func (acfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if acfg.Platform == "dev" {
		log.Printf("Deleting all users from DB...")
		_, err := acfg.DatabaseQueries.ResetUsers(context.Background())
		if err != nil {
			log.Printf("Error resetting users table in database: %v", err)
		}

		w.WriteHeader(200)
		acfg.fileserverHits.Store(0)
		log.Printf("Reset hits counter to 0.")
		w.Header().Add("Content-Type", "text/plan; charset=utf-8")
		w.Write([]byte("Reset hits counter!"))
	} else {
		log.Printf("Prevented unsafe ResetUsers operation on database, as not running in a dev environment")
		w.WriteHeader(403)
	}
}
