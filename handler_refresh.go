package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Rhyster42/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerValidateRefreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve bearer token", err)
		return
	}
	if refreshToken == "" {
		respondWithError(w, http.StatusUnauthorized, "token doesn't exist", errors.New("token doesn't exist"))
		return
	}

	databaseToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error retrieving Refresh Token from database", err)
		return
	}

	if time.Now().After(databaseToken.ExpiresAt.Time) || databaseToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token doesn't exist", err)
		return
	}

	newToken, err := auth.MakeJWT(databaseToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new Access Token", err)
		return
	}

	jsonParams := Token{
		Token: newToken,
	}

	respondWithJSON(w, http.StatusOK, jsonParams)
}
