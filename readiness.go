package main

import (
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	//sets the content-header.
	w.Header().Set("Content-Type", "text/html")

	//displays the status code.
	w.WriteHeader(200)

	//displays the body text.
	message := "OK"
	w.Write([]byte(message))
}
