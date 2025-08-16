package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Xeninon/Chirpy/internal/auth"
	"github.com/Xeninon/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, req *http.Request) {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResponse(w, "Unauthorized", 401)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		errorResponse(w, "Unauthorized", 401)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	if len(params.Body) >= 140 {
		errorResponse(w, "Chirp is too long", 400)
		return
	}

	params.Body = cleanBody(params.Body)
	chirp, err := cfg.db.CreateChirp(
		req.Context(),
		database.CreateChirpParams{
			Body:   params.Body,
			UserID: userID,
		},
	)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	dat, err := json.Marshal(
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    userID,
		},
	)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}

func cleanBody(body string) string {
	profanes := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		if slices.Contains(profanes, lowered) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, req *http.Request) {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	Chirps := make([]Chirp, 0)

	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	for _, chirp := range chirps {
		Chirps = append(
			Chirps,
			Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			},
		)
	}

	dat, err := json.Marshal(Chirps)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, req *http.Request) {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	userID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		errorResponse(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), userID)
	if err != nil {
		errorResponse(w, "Content not found", http.StatusNotFound)
		return
	}

	dat, err := json.Marshal(
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
