package database

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/keshavsharma98/chirpy/internal/common"
)

func (db *DB) RevokeRefreshToken(refresh_token string) error {
	refresh_token = strings.TrimPrefix(refresh_token, "Bearer ")
	dbData, err := db.readFromFile()
	if err != nil {
		return err
	}

	dbData.RevokedTokens[refresh_token] = time.Now()

	marshData, err := json.Marshal(dbData)
	if err != nil {
		return err
	}

	err = db.writeToFile(marshData)
	if err != nil {
		log.Println("Error while writing to data")
		return err
	}

	return nil
}

func (db *DB) RefreshToken(key, refresh_token string) (RefreshTokenResponse, error) {
	dbData, err := db.readFromFile()
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	id, err := common.CheckValidRefreshToken(key, refresh_token, dbData.RevokedTokens)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	token, err := common.CreateJWTToken(id, "chirpy-access", key)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	dbData.RevokedTokens[refresh_token] = time.Now()

	marshData, err := json.Marshal(dbData)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	err = db.writeToFile(marshData)
	if err != nil {
		log.Println("Error while writing to data")
		return RefreshTokenResponse{}, err
	}

	return RefreshTokenResponse{
		Token: token,
	}, nil
}
