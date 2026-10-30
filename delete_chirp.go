package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/MansurEh99/chirpy/internal/auth"
)

func (cfg *apiConfig) ChirpDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token missing", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	chirpID := r.PathValue("chirpID")
	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", err)
		return
	}

	ctx := r.Context()
	dbChirp, err := cfg.dbQueries.GetChirp(ctx, parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid chirp ID", err)
		return
	}

	if id != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "Invalid user ID", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(ctx, dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
