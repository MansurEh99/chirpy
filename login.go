package main

import (
	"net/http"
	"encoding/json"

	"github.com/MansurEh99/chirpy/internal/auth"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	var pass Pass
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pass)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	ctx := r.Context()
	dbUser,err := cfg.dbQueries.GetUser(ctx, pass.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(pass.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if match == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var user User
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email

	respondWithJSON(w, http.StatusOK, user)
}
