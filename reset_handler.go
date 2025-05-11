package main

import (
	"log"
	"net/http"
)

func (acfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	acfg.fileserverHits.Store(0)
	log.Printf("Reset hits counter to 0.")
	w.Header().Add("Content-Type", "text/plan; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Reset hits counter!"))
}
