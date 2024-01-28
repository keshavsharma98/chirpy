package database

import (
	"encoding/json"
	"errors"
	"log"
)

func (db *DB) CreateChirp(data string) (Chirp, error) {
	chirpsData, err := db.readFromFile()
	if err != nil {
		return Chirp{}, err
	}

	id := len(chirpsData.Chirps) + 1
	chirp := Chirp{
		Id:   id,
		Body: data,
	}
	chirpsData.Chirps[id] = chirp

	marshData, err := json.Marshal(chirpsData)
	if err != nil {
		return Chirp{}, err
	}

	err2 := db.writeToFile(marshData)
	if err2 != nil {
		log.Println("Error while writing to data")
		return Chirp{}, err2
	}
	id++
	return chirp, nil
}

func (db *DB) GetAllChirps() ([]Chirp, error) {
	chirpsData, err := db.readFromFile()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, len(chirpsData.Chirps))
	for _, v := range chirpsData.Chirps {
		chirps[v.Id-1] = v
	}
	return chirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	chirpsData, err := db.readFromFile()
	if err != nil {
		return Chirp{}, err
	}
	c, ok := chirpsData.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp does not exist")
	}
	return c, nil
}
