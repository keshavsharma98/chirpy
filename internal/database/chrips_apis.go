package database

import (
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strconv"
)

func (db *DB) CreateChirp(author_id int, data string) (Chirp, error) {
	chirpsData, err := db.readFromFile()
	if err != nil {
		return Chirp{}, err
	}

	id := len(chirpsData.Chirps) + 1
	chirp := Chirp{
		Id:       id,
		Body:     data,
		AuthorID: author_id,
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

func (db *DB) GetAllChirps(author, order string) ([]Chirp, error) {

	chirpsData, err := db.readFromFile()
	if err != nil {
		return nil, err
	}

	author_id := 0
	if author != "" {
		author_id_int, err := strconv.Atoi(author)
		if err != nil {
			return []Chirp{}, err
		}
		author_id = author_id_int
	}

	chirps := make([]Chirp, 0, len(chirpsData.Chirps))
	for _, v := range chirpsData.Chirps {
		if author != "" {
			if v.AuthorID == author_id {
				chirps = append(chirps, v)
			}
		} else {
			chirps = append(chirps, v)
		}
	}

	isAsc := order != "desc"

	sort.Slice(chirps, func(i, j int) bool {
		if isAsc {
			return chirps[i].Id < chirps[j].Id
		}
		return chirps[i].Id > chirps[j].Id
	})
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

func (db *DB) DeleteChirpById(user_id, id int) error {
	chirpsData, err := db.readFromFile()
	if err != nil {
		return err
	}

	c, ok := chirpsData.Chirps[id]
	if !ok {
		return errors.New("notfound")
	}

	if c.AuthorID != user_id {
		return errors.New("forbidden")
	}

	delete(chirpsData.Chirps, id)

	marshData, err := json.Marshal(chirpsData)
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
