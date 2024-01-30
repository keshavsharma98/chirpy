package main

import (
	"net/http"
)

func (apiCfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	err := apiCfg.DB.RevokeRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while revoking token")
	}
	respondWithJSON(w, 200, "success")
}

func (apiCfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	payloadBody, err := apiCfg.DB.RefreshToken(apiCfg.jwtSecret, token)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error while refreshing token")
	}
	respondWithJSON(w, 200, payloadBody)
}
