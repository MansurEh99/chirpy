package main

import (
	"net/http"
	"database/sql"
	"errors"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) ChirpsGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := r.URL.Query().Get("author_id")
	if s != "" {
		parsedUserID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
			return
		}

		dbChirps, err := cfg.dbQueries.GetChirpFromUserID(ctx, parsedUserID)
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No chirps available", err)
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to get chirps", err)
			return
		}

		ord := r.URL.Query().Get("sort")
		if ord == "" || ord == "asc" {
			sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.Before(dbChirps[j].CreatedAt) })
		}
		if ord == "desc" {
			sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[j].CreatedAt.Before(dbChirps[i].CreatedAt)})
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
		return
	}

	dbChirps, err := cfg.dbQueries.GetChirps(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "No chirps available", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get chirps", err)
		return
	}

	ord := r.URL.Query().Get("sort")
	if ord == "" || ord == "asc" {
		sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.Before(dbChirps[j].CreatedAt) })
	}
	if ord == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[j].CreatedAt.Before(dbChirps[i].CreatedAt)})
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
