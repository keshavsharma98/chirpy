package common

import (
	"errors"
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

func CreateJWTToken(id int, issuer, key string) (string, error) {
	expires_in_seconds := time.Duration(60) * 24 * time.Hour
	if issuer == "chiper-access" {
		expires_in_seconds = time.Second * 1 * 60 * 60
	}

	current_time := time.Now()
	expire_time := time.Now().Add(expires_in_seconds)
	registered_claims := jwt.RegisteredClaims{
		Issuer:    issuer,
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
	token, err := jwt.ParseWithClaims(signed_token_string, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		log.Println("Error parsing token:", err)
		return 0, errors.New("unauthorized")
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		if claims.Issuer == "chirpy-refresh" {
			log.Println("Invalid Issuer")
			return 0, errors.New("unauthorized")
		}
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

func CheckValidRefreshToken(key, request_token_string string, revoked_tokens map[string]time.Time) (int, error) {
	signed_token_string := strings.TrimPrefix(request_token_string, "Bearer ")

	_, exist := revoked_tokens[signed_token_string]
	if exist {
		return 0, errors.New("unauthorized")
	}

	token, err := jwt.ParseWithClaims(signed_token_string, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		log.Println("Error parsing token:", err)
		return 0, errors.New("unauthorized")
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		if claims.Issuer == "chirpy-refresh" {
			id, err := strconv.Atoi(claims.Subject)
			if err != nil {
				return 0, err
			}
			return id, nil
		}
	} else {
		log.Println("unknown claims type, cannot proceed")
		return 0, errors.New("unauthorized")
	}
	return 0, nil
}
