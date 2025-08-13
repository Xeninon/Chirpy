package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func errorResponse(w http.ResponseWriter, msg string, code int) {
	type returnVals struct {
		Error string `json:"error"`
	}
	respBody := returnVals{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
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
