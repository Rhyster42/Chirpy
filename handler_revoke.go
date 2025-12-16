package main

import (
	"net/http"

	"github.com/Rhyster42/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retreive bearer token", err)
		return
	}

	err = cfg.db.UpdateRevokedAt(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
