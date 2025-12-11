package main

import (
	"encoding/json"
	"net/http"
	"strings"

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
	params.Body = handlerValidateChirp(w, params.Body)

	dbChirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.User_id,
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
