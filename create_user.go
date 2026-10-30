package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/MansurEh99/chirpy/internal/database"
	"github.com/MansurEh99/chirpy/internal/auth"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool 	`json:"is_chirpy_red"`
}

type Pass struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

func (cfg *apiConfig) UserCreate(w http.ResponseWriter, r *http.Request) {
	var pass Pass
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pass)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode request", err)
		return
	}

	hash, err := auth.HashPassword(pass.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	var params database.CreateUserParams
	params.Email = pass.Email
	params.HashedPassword = hash
	ctx := r.Context()

	dbUser, err := cfg.dbQueries.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("Database Error: %v\n", err)
		return
	}

	var user User
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email
	user.IsChirpyRed = dbUser.IsChirpyRed

	respondWithJSON(w, http.StatusCreated, user)
}
