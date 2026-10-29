package main

import (
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/MansurEh99/chirpy/internal/auth"
	"github.com/MansurEh99/chirpy/internal/database"
)

type LogInfo struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type LoggedUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	var info LogInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&info)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	ctx := r.Context()
	dbUser,err := cfg.dbQueries.GetUser(ctx, info.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(info.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if match == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	hour, err := time.ParseDuration("1h")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot parse hour", err)
		return
	}
	jwt, err := auth.MakeJWT(dbUser.ID, cfg.secretKey, hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	refToken := auth.MakeRefreshToken()
	sixtyDays, err := time.ParseDuration("1440h")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot parse hour", err)
		return
	}
	params := database.CreateRefreshTokenParams{}
	params.Token = refToken
	params.UserID = dbUser.ID
	params.ExpiresAt = time.Now().Add(sixtyDays)
	dbRefToken, err := cfg.dbQueries.CreateRefreshToken(ctx, params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create refresh token", err)
		return
	}

	var user LoggedUser
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email
	user.Token = jwt
	user.RefreshToken = dbRefToken.Token

	respondWithJSON(w, http.StatusOK, user)
}
