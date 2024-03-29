package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/keshavsharma98/chirpy/internal/common"
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
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	body := reqBody{}

	err := decoder.Decode(&body)
	if err != nil {
		log.Println("Something went wrong", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	payloadBody, err := apiCfg.DB.Login(apiCfg.jwtSecret, body.Email, body.Password)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while login")
		return
	}
	respondWithJSON(w, 200, payloadBody)
}

func (apiCfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	id, err := common.CheckAuthorization(apiCfg.jwtSecret, token)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while creating chirp")
	}

	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	body := reqBody{}

	err = decoder.Decode(&body)
	if err != nil {
		log.Println("Something went wrong", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	payloadBody, err := apiCfg.DB.UpdateUser(id, body.Email, body.Password)
	if err != nil {
		if err.Error() == "user does not exist" {
			respondWithError(w, http.StatusNotFound, "user does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while updating user")
	}
	respondWithJSON(w, 200, payloadBody)
}

func (apiCfg *apiConfig) handlerWebhookUpgradeUser(w http.ResponseWriter, r *http.Request) {
	polka_key := os.Getenv("POLKA_KEY")
	api_key := r.Header.Get("Authorization")
	api_key = strings.TrimPrefix(api_key, "ApiKey ")

	if polka_key != api_key {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	type reqBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	body := reqBody{}

	err := decoder.Decode(&body)
	if err != nil {
		log.Println("Something went wrong", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if body.Event != "user.upgraded" {
		respondWithJSON(w, 200, nil)
		return
	}

	payloadBody, err := apiCfg.DB.WebhookUpgradeUser(body.Data.UserID)
	if err != nil {
		if err.Error() == "user does not exist" {
			respondWithError(w, http.StatusNotFound, "user does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while upgrading user")
	}
	respondWithJSON(w, 200, payloadBody)
}
