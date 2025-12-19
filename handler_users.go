package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Rhyster42/Chirpy/internal/auth"
	"github.com/Rhyster42/Chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type upgradeParams struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	}
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
		ID:          data.ID,
		CreatedAt:   data.CreatedAt.Time,
		UpdatedAt:   data.UpdatedAt.Time,
		Email:       data.Email,
		IsChirpyRed: false,
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
		IsChirpyRed:  user.IsChirpyRed.Bool,
	}
	respondWithJSON(w, http.StatusOK, returnedUser)

}

func (cfg *apiConfig) handlerUpdateEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	//pulling/veriying access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to retrieve Access Token", err)
		return
	}
	if accessToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid/missing token", err)
		return
	}

	// retrieve New Email and Password from Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access Token", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	newUser, err := cfg.db.ChangeEmailAndPassword(r.Context(),
		database.ChangeEmailAndPasswordParams{
			ID:             userID,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to change Email and Password.", err)
		return
	}

	userResponse := User{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt.Time,
		UpdatedAt:   newUser.UpdatedAt.Time,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {

	requestAPI, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid/missing API header", err)
		return
	}
	if requestAPI != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Incorrect API key", errors.New("Incorrect API"))
	}

	decoder := json.NewDecoder(r.Body)
	params := upgradeParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to find user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
