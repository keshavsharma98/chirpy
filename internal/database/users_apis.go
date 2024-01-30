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

func (db *DB) Login(jwt_secret_key, email, password string, expires_in_seconds int) (UserResponse, error) {
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

	token, err := common.CreateJWTToken(expires_in_seconds, user.Id, jwt_secret_key)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: token,
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
