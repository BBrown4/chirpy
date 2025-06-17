package main

import (
	"encoding/json"
	"net/http"

	"github.com/bbrown4/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		},
	})
}
