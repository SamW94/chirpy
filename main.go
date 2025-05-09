package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (acfg *apiConfig) middlewareFileServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incrementing counter for: %s", r.URL.Path)
		acfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (acfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	hitsCount := acfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %v", hitsCount)))
}

func (acfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	acfg.fileserverHits.Store(0)
	log.Printf("Reset hits counter to 0.")
	w.Header().Add("Content-Type", "text/plan; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Reset hits counter!"))
}

func main() {
	const serverPort = "8080"
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /reset", apiCfg.resetHandler)

	appPathHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app", apiCfg.middlewareFileServerHits(appPathHandler))
	mux.Handle("/app/", apiCfg.middlewareFileServerHits(appPathHandler))

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	log.Printf("Serving on port: %s", serverPort)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
