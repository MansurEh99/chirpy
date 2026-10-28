package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Email struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) UserCreate(w http.ResponseWriter, r *http.Request) {
	var email Email
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&email)
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return
	}

	ctx := r.Context()
	dbUser, err := cfg.dbQueries.CreateUser(ctx, email.Email)
	if err != nil {
		fmt.Printf("Database Error: %v\n", err)
		return
	}

	var user User
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email

	respondWithJSON(w, http.StatusCreated, user)
}
