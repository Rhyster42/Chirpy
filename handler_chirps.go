package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Rhyster42/Chirpy/internal/auth"
	"github.com/Rhyster42/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Response struct {
	Error string `json:"error"`
	Valid bool   `json:"valid"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request to create Chirp", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve token for user", err)
		return
	}

	authUser, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid/Expired token", err)
		return
	}

	params.Body = handlerValidateChirp(w, params.Body)

	dbChirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: authUser,
	}

	data, err := cfg.db.CreateChirp(r.Context(), dbChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp", err)
		return
	}

	chirp := Chirp{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserID:    data.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func handlerValidateChirp(w http.ResponseWriter, body string) string {

	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return ""
	}
	return checkProfanity(body)
}

func checkProfanity(body string) string {

	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	word_list := strings.Split(body, " ")

	for i, word := range word_list {
		lowercase_word := strings.ToLower(word)
		for _, badword := range badWords {
			if lowercase_word == badword {
				word_list[i] = "****"
			}
		}
	}
	return strings.Join(word_list, " ")
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to retrieve Access Token", err)
		return
	}
	if accessToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid/missing token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access Token", err)
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Must be author of chirp to delete it", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
