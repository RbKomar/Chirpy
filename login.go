package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RbKomar/Chirpy/internal/auth"
)

type UserJWT struct {
	User
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parametres", err)
		return
	}

	var expirationTime time.Duration
	if params.ExpiresInSeconds == 0 {
		expirationTime = time.Hour
	} else if params.ExpiresInSeconds > 3600 {
		expirationTime = time.Hour
	} else {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	user, err := cfg.dbQueries.PasswordLookupByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't find user with matching email", err)
		return
	}
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Wrong password match", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expirationTime)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}
	RespondWithJson(w, http.StatusOK, UserJWT{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			Email:     user.Email,
		},
		Token: token,
	})
}
