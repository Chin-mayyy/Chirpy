package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(404)
		w.Write([]byte("Forbidden"))
	} else {
		w.WriteHeader(200)
		cfg.db.Reset(r.Context())
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits resets to 0"))
}
