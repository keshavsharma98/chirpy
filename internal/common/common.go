package common

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) string {
	enc_pswrd, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		log.Fatalln("Failed at encrypting password", err)
		return ""
	}
	return string(enc_pswrd)
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
