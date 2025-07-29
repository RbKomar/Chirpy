package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RbKomar/Chirpy/internal/auth"
	"github.com/RbKomar/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
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
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	dbParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), dbParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't add user to db", err)
		return
	}
	userJson := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	}
	RespondWithJson(w, http.StatusCreated, userJson)
}
