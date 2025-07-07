package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		RespondWithError(w, 400, "message is too long", nil)
		return
	}
	HandleBannedWords(w, params.Body)
}

func HandleBannedWords(w http.ResponseWriter, message string) {
	type returnValid struct {
		CleanedBody string `json:"cleaned_body"`
	}

	sep := " "
	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(message, sep)
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := bannedWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	censoredMessage := strings.Join(words, sep)
	RespondWithJson(w, http.StatusOK, returnValid{
		CleanedBody: censoredMessage,
	})
}
