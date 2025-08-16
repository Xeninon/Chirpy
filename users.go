package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Xeninon/Chirpy/internal/auth"
	"github.com/Xeninon/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	type parameters struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresIn int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenDuration := min(time.Duration(params.ExpiresIn)*time.Second, time.Hour)
	if params.ExpiresIn == 0 {
		tokenDuration = time.Hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, tokenDuration)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	dat, err := json.Marshal(
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
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

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, req *http.Request) {
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	user, err := cfg.db.CreateUser(
		req.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hash,
		},
	)
	if err != nil {
		errorResponse(w, "Something went wrong", 500)
		return
	}

	dat, err := json.Marshal(
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
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
