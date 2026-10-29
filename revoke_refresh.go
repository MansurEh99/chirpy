package main

import (
	"net/http"

	"github.com/MansurEh99/chirpy/internal/auth"
)

func (cfg *apiConfig) RefreshRevoke (w http.ResponseWriter, r *http.Request) {
	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token detected", err)
		return
	}

	ctx := r.Context()
	err = cfg.dbQueries.UpdateRefreshTokenRevoke(ctx, refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to update database", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
