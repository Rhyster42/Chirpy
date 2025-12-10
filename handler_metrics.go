package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerServerHits(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	msg := fmt.Sprintf(
		`<html>
			<body>
    			<h1>Welcome, Chirpy Admin</h1>
    			<p>Chirpy has been visited %d times!</p>
  			</body>
		</html>`, cfg.fileserverHits.Load())

	w.Write([]byte(msg))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
