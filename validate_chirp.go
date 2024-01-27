package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type reqBody struct {
	Body string `json:"body"`
}

var profaneWords = map[string]string{
	"kerfuffle": "****",
	"sharbert":  "****",
	"fornax":    "****",
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := reqBody{}

	err := decoder.Decode(&body)
	if err != nil {
		log.Println("Something went wrong", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanString := removeProfaneWords(body.Body)

	payloadBody := map[string]string{
		"cleaned_body": cleanString,
	}

	respondWithJSON(w, 200, payloadBody)
}

func removeProfaneWords(s string) string {
	s_arr := strings.Split(s, " ")

	for i, e := range s_arr {
		v, ok := profaneWords[strings.ToLower(e)]
		if ok {
			s_arr[i] = v
		}
	}
	return strings.Join(s_arr, " ")
}
