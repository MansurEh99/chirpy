package main

import (
	"net/http"
	"errors"
	"encoding/json"
	"database/sql"

	"github.com/google/uuid"
	"github.com/MansurEh99/chirpy/internal/auth"
)

type UpgradeInfo struct {
	Event string `json:"event"`
	Data struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) UserUpgrade(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}

	var info UpgradeInfo
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&info)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if info.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	parsedUserID, err := uuid.Parse(info.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server unable to parse ID", err)
		return
	}
	ctx := r.Context()
	_, err = cfg.dbQueries.UpgradeUser(ctx, parsedUserID)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "Unregistered user", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to upgrade", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
