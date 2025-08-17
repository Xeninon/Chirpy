package main

import (
	"encoding/json"
	"net/http"

	"github.com/Xeninon/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type PolkaRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) polkaHandler(w http.ResponseWriter, req *http.Request) {
	polkaKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if polkaKey != cfg.polkaKey {
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := PolkaRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if err = cfg.db.UpgradeToRed(req.Context(), userID); err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}
