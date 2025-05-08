package main

import (
	"log"
	"net/http"

	"github.com/SamW94/chirpy/internal/http_handlers"
)

func main() {
	const serverPort = "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/healthz", http_handlers.HealthzHandler)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	log.Printf("Serving on port: %s", serverPort)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
