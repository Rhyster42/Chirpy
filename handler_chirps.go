package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Response struct {
	Error string `json:"error"`
	Valid bool   `json:"valid"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
		Valid        bool   `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_Body: checkProfanity(params.Body),
		Valid:        true,
	})
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
