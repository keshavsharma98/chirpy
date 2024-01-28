package database

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/keshavsharma98/chirpy/internal/common"
)

func (db *DB) CreateUsers(email, password string) (UserResponse, error) {
	dbData, err := db.readFromFile()
	if err != nil {
		return UserResponse{}, err
	}

	for _, v := range dbData.Users {
		if v.Email == email {
			return UserResponse{}, errors.New("username already exists")
		}
	}

	id := len(dbData.Users) + 1
	hash := common.EncryptPassword(password)
	user := User{
		Id:       id,
		Email:    email,
		Password: hash,
	}
	dbData.Users[id] = user

	marshData, err := json.Marshal(dbData)
	if err != nil {
		return UserResponse{}, err
	}

	err2 := db.writeToFile(marshData)
	if err2 != nil {
		log.Println("Error while writing to data")
		return UserResponse{}, err2
	}
	return UserResponse{
		Email: user.Email,
		Id:    user.Id,
	}, nil
}

func (db *DB) Login(email, password string) (UserResponse, error) {
	dbData, err := db.readFromFile()
	if err != nil {
		return UserResponse{}, err
	}

	user := User{}
	isExist := false

	for _, v := range dbData.Users {
		if v.Email == email {
			user = v
			isExist = !isExist
			break
		}
	}
	if !isExist {
		return UserResponse{}, errors.New("unauthorized")
	}

	err = common.ComparePassword(user.Password, password)
	if err != nil {
		return UserResponse{}, errors.New("unauthorized")
	}

	return UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}, nil
}
