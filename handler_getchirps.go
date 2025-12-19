package main

import (
	"net/http"
	"slices"

	"github.com/Rhyster42/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	authorString := r.URL.Query().Get("author_id")

	var data []database.Chirp
	var err error
	if authorString == "" {
		data, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
			return
		}
	} else {
		authorID, err := uuid.Parse(authorString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse ID", err)
			return
		}
		data, err = cfg.db.GetChirpsFromUser(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve Authors Chirps", err)
			return
		}
	}
	sort := r.URL.Query().Get("sort")

	chirps := make([]Chirp, len(data))

	for i, item := range data {
		chirps[i] = Chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID:    item.UserID,
		}
	}

	if sort == "desc" {
		slices.Reverse(chirps)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpIDString := r.PathValue("chirpID")
	u, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID", err)
	}

	data, err := cfg.db.GetChirp(r.Context(), u)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp in database", err)
	}
	chirp := Chirp{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserID:    data.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)

}
