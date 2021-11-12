package token

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/keruch/tfs-go-hw/hw4/internal/domain"
)

var ErrorKeyNotFound = errors.New("encryption/decryption  key not found")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateUserToken(data domain.UserData) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: data.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtKey := []byte(os.Getenv("ACCESS_SECRET"))
	if jwtKey == nil {
		return "", ErrorKeyNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateUserToken(tokenString string) (string, error) {
	claims := &Claims{}
	jwtKey := []byte(os.Getenv("ACCESS_SECRET"))
	if jwtKey == nil {
		return "", ErrorKeyNotFound
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", nil
	}

	return claims.Username, nil
}
