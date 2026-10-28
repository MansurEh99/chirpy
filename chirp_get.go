package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) ChirpGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error parsing chirp ID", err)
		return
	}

	ctx := r.Context()
	dbChirp, err := cfg.dbQueries.GetChirp(ctx, parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp", err)
		return
	}

	var chirp Chirp
	chirp.ID = dbChirp.ID
	chirp.CreatedAt = dbChirp.CreatedAt
	chirp.UpdatedAt = dbChirp.UpdatedAt
	chirp.Body = dbChirp.Body
	chirp.UserID = dbChirp.UserID

	respondWithJSON(w, http.StatusOK, chirp)
}
