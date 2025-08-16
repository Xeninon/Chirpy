package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Xeninon/Chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(
		req.Context(),
		token,
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(user.UserID, cfg.secret, time.Hour)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	dat, err := json.Marshal(
		Response{
			Token: accessToken,
		},
	)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
