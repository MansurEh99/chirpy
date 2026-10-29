package main

import (
	"strings"
	"net/http"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/MansurEh99/chirpy/internal/database"
	"github.com/MansurEh99/chirpy/internal/auth"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID 	  uuid.UUID `json:"user_id"`
}

type ChirpParams struct {
	Body   string 		`json:"body"`
	UserID uuid.UUID 	`json:"user_id"`
}

func (cfg *apiConfig) ChirpCreate (w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params ChirpParams
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxLength = 140
	if len(params.Body) > maxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not logged in", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	var chirpParams database.CreateChirpParams
	chirpParams.Body = params.Body
	chirpParams.UserID = id
	ctx := r.Context()

	dbChirp, err := cfg.dbQueries.CreateChirp(ctx, chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp", err)
		return
	}

	var chirp Chirp
	chirp.ID = dbChirp.ID
	chirp.CreatedAt = dbChirp.CreatedAt
	chirp.UpdatedAt = dbChirp.UpdatedAt
	chirp.Body = dbChirp.Body
	chirp.UserID = dbChirp.UserID

	respondWithJSON(w, http.StatusCreated, chirp)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: msg,})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func cleaned(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
