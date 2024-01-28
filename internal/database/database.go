package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DatabaseSchema struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	i, _ := os.Stat("database.json")
	if i != nil {
		log.Println("database exists")
		return db, nil
	}

	log.Println("Cannot find DB. Will create new DB")

	schema := DatabaseSchema{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	err2 := db.createDB(&schema)
	if err2 != nil {
		log.Println("Error while creating new DB")
		return nil, err2
	}
	return db, nil
}

func (db *DB) createDB(schema *DatabaseSchema) error {
	data, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	err2 := db.writeToFile(data)
	if err2 != nil {
		return err2
	}
	return nil
}

func (db *DB) readFromFile() (DatabaseSchema, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := DatabaseSchema{}
	f, err := os.ReadFile(db.path)
	if err != nil {
		log.Println("Error while reading file")
		return DatabaseSchema{}, err
	}

	err = json.Unmarshal(f, &data)
	if err != nil {
		return DatabaseSchema{}, err
	}

	return data, nil
}

func (db *DB) writeToFile(data []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	err := os.WriteFile(db.path, data, 0600)
	if err != nil {
		log.Println("Error while writing to file")
		return err
	}
	return nil
}
