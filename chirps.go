package main

import (
	"net/http"
)

func (cfg *apiConfig) ChirpsGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dbChirps, err := cfg.dbQueries.GetChirps(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps from database", err)
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i].ID = dbChirp.ID
		chirps[i].CreatedAt = dbChirp.CreatedAt
		chirps[i].UpdatedAt = dbChirp.UpdatedAt
		chirps[i].Body = dbChirp.Body
		chirps[i].UserID = dbChirp.UserID
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
