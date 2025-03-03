package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	type errorMessage struct {
		Error string `json:"error"`
	}

	if code >= 400 && code < 500 {
		mess := errorMessage{
			Error: message,
		}

		data, err := json.Marshal(mess)
		if err != nil {
			log.Fatalf("error marshaling the error message : %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	}
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func replaceProfane(body string) string {
	newBody := []string{}

	oldBody := strings.Split(body, " ")
	for _, word := range oldBody {
		switch strings.ToLower(word) {
		case "kerfuffle":
			newBody = append(newBody, "****")
		case "sharbert":
			newBody = append(newBody, "****")
		case "fornax":
			newBody = append(newBody, "****")
		default:
			newBody = append(newBody, word)
		}
	}

	return strings.Join(newBody, " ")
}
