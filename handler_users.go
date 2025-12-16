package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rhyster42/Chirpy/internal/auth"
	"github.com/Rhyster42/Chirpy/internal/database"
)

type parameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedpassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
	}
	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedpassword,
	}

	data, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
	}

	user := User{
		ID:        data.ID,
		CreatedAt: data.CreatedAt.Time,
		UpdatedAt: data.UpdatedAt.Time,
		Email:     data.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode JSON", err)
		return
	}

	//Password
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving user from database", err)
	}

	valid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error with password authentication", err)
	}
	if !valid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password", err)
	}

	//JWT
	expireTime := time.Duration(time.Hour)

	token, err := auth.MakeJWT(user.ID, cfg.secret, expireTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retreive token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add refresh token to database.", err)
		return
	}

	returnedUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, returnedUser)

}
