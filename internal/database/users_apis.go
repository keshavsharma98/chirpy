package database

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/keshavsharma98/chirpy/internal/common"
)

func (db *DB) CreateUsers(email, password string) (UpdateUserResponse, error) {
	dbData, err := db.readFromFile()
	if err != nil {
		return UpdateUserResponse{}, err
	}

	for _, v := range dbData.Users {
		if v.Email == email {
			return UpdateUserResponse{}, errors.New("username already exists")
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
		return UpdateUserResponse{}, err
	}

	err2 := db.writeToFile(marshData)
	if err2 != nil {
		log.Println("Error while writing to data")
		return UpdateUserResponse{}, err2
	}
	return UpdateUserResponse{
		Email: user.Email,
		Id:    user.Id,
	}, nil
}

func (db *DB) Login(jwt_secret_key, email, password string) (UserLoginResponse, error) {
	dbData, err := db.readFromFile()
	if err != nil {
		return UserLoginResponse{}, err
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
		return UserLoginResponse{}, errors.New("unauthorized")
	}

	err = common.ComparePassword(user.Password, password)
	if err != nil {
		return UserLoginResponse{}, errors.New("unauthorized")
	}

	token, err := common.CreateJWTToken(user.Id, "chirpy-access", jwt_secret_key)
	if err != nil {
		return UserLoginResponse{}, err
	}

	refresh_token, err := common.CreateJWTToken(user.Id, "chirpy-refresh", jwt_secret_key)
	if err != nil {
		return UserLoginResponse{}, err
	}

	return UserLoginResponse{
		Id:           user.Id,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token,
	}, nil
}

func (db *DB) UpdateUser(id int, email, password string) (UpdateUserResponse, error) {
	dbData, err := db.readFromFile()
	if err != nil {
		return UpdateUserResponse{}, err
	}

	user := User{}

	isUserExist := false
Loop_Update_User:
	for _, v := range dbData.Users {
		if v.Id == id {
			user = v
			isUserExist = !isUserExist
			break Loop_Update_User
		}
	}

	if !isUserExist {
		return UpdateUserResponse{}, errors.New("user does not exist")
	}

	hash := common.EncryptPassword(password)
	user.Email = email
	user.Password = hash
	dbData.Users[id] = user

	marshData, err := json.Marshal(dbData)
	if err != nil {
		return UpdateUserResponse{}, err
	}

	err2 := db.writeToFile(marshData)
	if err2 != nil {
		log.Println("Error while writing to data")
		return UpdateUserResponse{}, err2
	}
	return UpdateUserResponse{
		Email: user.Email,
		Id:    user.Id,
	}, nil
}
