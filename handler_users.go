package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	body := reqBody{}

	err := decoder.Decode(&body)
	if err != nil {
		log.Println("Something went wrong", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	payloadBody, err := apiCfg.DB.CreateUsers(body.Email, body.Password)
	if err != nil {
		if err.Error() == "username already exists" {
			respondWithError(w, http.StatusUnauthorized, "username already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while creating user")
	}
	respondWithJSON(w, 201, payloadBody)
}

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	body := reqBody{}

	err := decoder.Decode(&body)
	if err != nil {
		log.Println("Something went wrong", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	payloadBody, err := apiCfg.DB.Login(body.Email, body.Password)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while creating chirp")
	}
	respondWithJSON(w, 200, payloadBody)
}
