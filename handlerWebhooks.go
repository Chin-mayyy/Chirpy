package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Chin-mayyy/Chirpy/internal/auth"
	"github.com/Chin-mayyy/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type Parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	api_key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting pai ke", err)
	}
	if api_key != cfg.polkaKey {
		respondWithJSON(w, http.StatusUnauthorized, "Api key is not valid")
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 404, "Error decoding the JSON Body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
	}

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 400, "Invalid user id format", err)
	}

	_, err = cfg.db.UpdateUserChirpy(r.Context(), database.UpdateUserChirpyParams{
		ID:          id,
		IsChirpyRed: sql.NullBool{Bool: true, Valid: true},
	})

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error updating user: %v", err), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
