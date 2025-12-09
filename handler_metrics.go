package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerServerHits(w http.ResponseWriter, r *http.Request) {

	msg := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())

	w.Write([]byte(msg))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}
