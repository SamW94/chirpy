package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SamW94/chirpy/internal/auth"
	"github.com/SamW94/chirpy/internal/database"
	"github.com/google/uuid"
)

type requestJSONCreateUser struct {
	Password      string `json:"password"`
	RequestedUser string `json:"email"`
}

type CreateUserResponseSuccessful struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (acfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestJson := requestJSONCreateUser{}
	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Printf("Error decoding request body to requestJSON: %v", err)
		respondWithError(w, 500, "Something went wrong.")
		return
	}

	if requestJson.Password == "" {
		log.Printf("No password provided in request body.")
		respondWithError(w, 500, "No password provided in request body.")
		return
	} else {
		hashedPassword, err := auth.HashPassword(requestJson.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			respondWithError(w, 500, "Something went wrong.")
			return
		}

		createUserParams := database.CreateUserParams{
			Email:          requestJson.RequestedUser,
			HashedPassword: hashedPassword,
		}

		dbUser, err := acfg.DatabaseQueries.CreateUser(context.Background(), createUserParams)
		if err != nil {
			log.Printf("Error calling database.CreateUser() function: %v", err)
			respondWithError(w, 500, "Something went wrong.")
			return
		}

		log.Printf("Successfully created user with ID %v and email %v", dbUser.ID, dbUser.Email)
		respBody := CreateUserResponseSuccessful{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		}

		respondWithJSON(w, 201, respBody)
	}
}
