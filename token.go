package main

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	Guid string `json:"guid"`
	jwt.StandardClaims
}

func createTokens(guid string) (string, string, error) {
	secret := []byte("secret")
	claims := Claims{
		guid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(secret)
	lenToken := len(tokenString)

	refresh := time.Now().Add(time.Hour * 24 * 60).UTC().String()
	refresh += tokenString[lenToken-4 : lenToken-1]
	refresh = base64.StdEncoding.EncodeToString([]byte(refresh))
	bytesRefresh, err := bcrypt.GenerateFromPassword([]byte(refresh), 14)
	createRefreshToken(bytesRefresh, guid)
	return tokenString, refresh, err
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
