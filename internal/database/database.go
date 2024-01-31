package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DatabaseSchema struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RevokedTokens map[string]time.Time `json:"revoked_tokens"`
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserLoginResponse struct {
	Email        string `json:"email"`
	Id           int    `json:"id"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateUserResponse struct {
	Email       string `json:"email"`
	Id          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
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
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RevokedTokens: map[string]time.Time{},
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
