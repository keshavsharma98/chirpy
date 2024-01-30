package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/keshavsharma98/chirpy/internal/common"
)

var profaneWords = map[string]string{
	"kerfuffle": "****",
	"sharbert":  "****",
	"fornax":    "****",
}

func (apiCfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
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
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	body := reqBody{}

	err = decoder.Decode(&body)
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

	payloadBody, err := apiCfg.DB.CreateChirp(id, cleanString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating chirp")
	}
	respondWithJSON(w, 201, payloadBody)
}

func (apiCfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	payloadBody, err := apiCfg.DB.GetAllChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}
	respondWithJSON(w, 200, payloadBody)
}

func (apiCfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chirp_id, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}

	payloadBody, err2 := apiCfg.DB.GetChirpById(chirp_id)
	if err2 != nil {
		respondWithError(w, http.StatusNotFound, "chirp does not exist")
		return
	}
	respondWithJSON(w, 200, payloadBody)
}

func (apiCfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	user_id, err := common.CheckAuthorization(apiCfg.jwtSecret, token)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while creating chirp")
	}

	id := chi.URLParam(r, "id")
	chirp_id, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}

	err = apiCfg.DB.DeleteChirpById(user_id, chirp_id)
	if err != nil {
		if err.Error() == "forbidden" {
			respondWithError(w, http.StatusForbidden, "Forbidden")
			return
		}
		if err.Error() == "notfound" {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while deleting chirp")
		return
	}
	respondWithJSON(w, 200, "deleted successfully")
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
