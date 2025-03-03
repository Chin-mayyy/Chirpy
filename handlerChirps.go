package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Chin-mayyy/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameter{}

	//Decoding the json.
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
	}

	//Checking the length of chirps
	if len(params.Body) >= 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		newBody := replaceProfane(params.Body)

		chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   newBody,
			UserID: params.UserID,
		})
		if err != nil {
			respondWithError(w, 400, "Error creating a chirp")
		}

		resp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			Body:      chirp.Body,
			UserID:    params.UserID.String(),
		}

		respondWithJson(w, 201, resp)
	}
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Error getting chirps")
	}

	resp := []Chirp{}

	for _, chirp := range chirps {
		resp = append(resp, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		})
	}

	respondWithJson(w, 200, resp)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	id := strings.TrimPrefix(url, "/api/chirps/")
	ID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Error parsing ID")
	}

	chirp, err := cfg.db.GetChirp(r.Context(), ID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
	}

	resp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}

	respondWithJson(w, 200, resp)
}
