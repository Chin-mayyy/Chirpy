package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Chin-mayyy/Chirpy/internal/auth"
	"github.com/Chin-mayyy/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
}

type Token struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerAcceptEmail(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameter{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add user")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	err = cfg.db.AddPassword(r.Context(), database.AddPasswordParams{
		HashedPassword: hashedPassword,
		Email:          user.Email,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add password")
	}

	respUser := User{
		ID:        user.ID,
		CreatedAt: user.UpdatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJson(w, 201, respUser)

	w.WriteHeader(201)
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := parameter{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Couldn't get user")
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, 401, "Invalid password")
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTsecret)
	if err != nil {
		respondWithError(w, 401, "Couldn't make the Token")
	}

	tokenString, _ := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     tokenString,
		CreatedAt: time.Now().UTC(),
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	respUser := User{
		ID:            user.ID,
		CreatedAt:     user.UpdatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Token:         token,
		Refresh_token: tokenString,
	}

	respondWithJson(w, 200, respUser)

	w.WriteHeader(200)
}

func (cfg *apiConfig) handlerRefreshes(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error getting refresh Token from header")
		return
	}

	if tokenString == "" {
		respondWithError(w, 401, "Refresh token doesn't exists")
		return
	}

	rToken, err := cfg.db.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, 401, "Error getting refresh token from database")
		return
	}

	if rToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Refresh token has expired")
		return
	}

	if rToken.RevokedAt.Valid {
		respondWithError(w, 401, "Refresh token has been revoked")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), rToken.UserID)
	if err != nil {
		respondWithError(w, 401, "Error getting user from database")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JWTsecret)
	if err != nil {
		respondWithError(w, 401, "Error creating access token")
		return
	}

	token := Token{
		Token: accessToken,
	}

	respondWithJson(w, 200, token)
}

func (cfg *apiConfig) handlerRevokes(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error getting refresh Token from header")
		return
	}

	if tokenString == "" {
		respondWithError(w, 401, "Refresh token doesn't exists")
		return
	}

	rToken, err := cfg.db.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, 401, "Error getting refresh token from database")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), rToken.Token)
	if err != nil {
		respondWithError(w, 401, "Error revoking refresh token from database")
		return
	}

	w.WriteHeader(204)
}
