package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/SamW94/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	const serverPort = "8080"
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits:  atomic.Int32{},
		DatabaseQueries: dbQueries,
	}

	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	appPathHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareFileServerHits(appPathHandler))

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	log.Printf("Serving on port: %s", serverPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
