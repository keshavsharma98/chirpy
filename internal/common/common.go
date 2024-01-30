package common

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func CreateJWTToken(expires_in_seconds, id int, key string) (string, error) {
	current_time := time.Now()
	expire_time := time.Now().Add(time.Second * time.Duration(expires_in_seconds))
	registered_claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(current_time),
		ExpiresAt: jwt.NewNumericDate(expire_time),
		Subject:   strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registered_claims)
	signed_token, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("Error in creating token: ", err)
		return "", err
	}
	return signed_token, nil
}

func CheckAuthorization(key, request_token_string string) (int, error) {
	signed_token_string := strings.TrimPrefix(request_token_string, "Bearer ")
	fmt.Println(signed_token_string)
	token, err := jwt.ParseWithClaims(signed_token_string, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		log.Println("Error parsing token:", err)
		return 0, errors.New("unauthorized")
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			log.Println("Error in token decoding", err)
			return 0, err
		}
		return id, nil
	} else {
		log.Println("unknown claims type, cannot proceed")
		return 0, errors.New("unauthorized")
	}
}
