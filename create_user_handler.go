package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type requestJSONCreateUser struct {
	RequestedUser string `json:"email"`
}

type CreateUserResponseSuccessful struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (acfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := requestJSONCreateUser{}
	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Printf("Error decoding request body to requestJSON: %v", err)
		respondWithError(w, 500, "Something went wrong")
	}

	dbQueries := acfg.DatabaseQueries
	dbUser, err := dbQueries.CreateUser(context.Background(), requestJson.RequestedUser)
	if err != nil {
		log.Printf("Error calling database.CreateUser() function: %v", err)
	}

	log.Printf("Successfully created user with ID %v and email %v", dbUser.ID, dbUser.Email)
	respBody := CreateUserResponseSuccessful{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, 201, respBody)

}
