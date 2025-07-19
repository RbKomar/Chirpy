package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/RbKomar/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func MapDBchirpToChirp(chirp database.Chirp) Chirp {
	return Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}

func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	chirps, err := cfg.dbQueries.RetrieveAllChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve all chirps", err)
		return
	}
	var jsonChirps []Chirp
	for _, chirp := range chirps {
		jsonChirps = append(jsonChirps, MapDBchirpToChirp(chirp))
	}

	RespondWithJson(w, http.StatusOK, jsonChirps)
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	if err := decoder.Decode(params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "can't decode params", err)
		return
	}

	censoredMessage, err := cfg.HandleValidation(w, params.Body)
	if err != nil {
		return
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   censoredMessage,
		UserID: params.UserID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't add chirp to DB.", err)
		return
	}
	chirpJson := MapDBchirpToChirp(chirp)
	RespondWithJson(w, http.StatusCreated, chirpJson)

}
func (cfg *apiConfig) HandleValidation(w http.ResponseWriter, message string) (string, error) {
	if len(message) > 140 {
		RespondWithError(w, 400, "message is too long", nil)
		return "", errors.New("message is too loong")
	}
	return HandleBannedWords(w, message), nil
}

func HandleBannedWords(w http.ResponseWriter, message string) string {
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
	return censoredMessage
}
