package main

import (
	"encoding/json"
	"net/http"
)

type errorResponseBody struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	body := errorResponseBody{
		Error: msg,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
