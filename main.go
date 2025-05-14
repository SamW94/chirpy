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
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	const serverPort = "8080"
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits:  atomic.Int32{},
		DatabaseQueries: dbQueries,
		Platform:        os.Getenv("PLATFORM"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		PolkaKey:        os.Getenv("POLKA_KEY"),
	}

	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", apiCfg.getChirpByIDHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirp_id}", apiCfg.deleteChirpByIDHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaWebHookHandler)

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
