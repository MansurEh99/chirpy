package main

import (
	"time"
	"net/http"

	"github.com/MansurEh99/chirpy/internal/auth"
)

type Refresh struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) RefreshGet (w http.ResponseWriter, r *http.Request) {
	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token detected", err)
		return
	}

	ctx := r.Context()
	dbUser, err := cfg.dbQueries.GetUserFromRefreshToken(ctx, refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to fetch user", err)
		return
	}

	var response Refresh
	jwt, err := auth.MakeJWT(dbUser.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not create access token", err)
		return
	}

	response.Token = jwt
	respondWithJSON(w, http.StatusOK, response)
}
