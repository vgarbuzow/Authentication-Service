package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	Uuid string `json:"uuid"`
	jwt.StandardClaims
}

func createTokens() (string, error) {
	claims := Claims{
		uuid.New().String(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	return tokenString, err
}

func parseAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if _, ok := token.Claims.(jwt.Claims); ok && token.Valid {
		return "All good", nil
	} else {
		return "", err
	}
}

func refreshTokens() {

}
