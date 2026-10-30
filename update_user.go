package main

import (
	"net/http"
	"encoding/json"

	"github.com/MansurEh99/chirpy/internal/database"
	"github.com/MansurEh99/chirpy/internal/auth"
)

func (cfg *apiConfig) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var info LogInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&info)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	hash, err := auth.HashPassword(info.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token missing", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	ctx := r.Context()
	var params database.UpdateUserParams
	params.Email = info.Email
	params.HashedPassword = hash
	params.ID = id
	updatedDbUser, err := cfg.dbQueries.UpdateUser(ctx, params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error updating database", err)
		return
	}

	var updatedUser User
	updatedUser.ID = updatedDbUser.ID
	updatedUser.CreatedAt = updatedDbUser.CreatedAt
	updatedUser.UpdatedAt = updatedDbUser.UpdatedAt
	updatedUser.Email = updatedDbUser.Email
	updatedUser.IsChirpyRed = updatedDbUser.IsChirpyRed

	respondWithJSON(w, http.StatusOK, updatedUser)
}
