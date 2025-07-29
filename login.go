package main

import (
	"encoding/json"
	"net/http"

	"github.com/RbKomar/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parametres", err)
		return
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
	RespondWithJson(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	})

}
