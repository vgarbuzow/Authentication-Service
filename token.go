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

func CreateTokens(guid string) (string, string, error) {
	secret := []byte("secret")
	claims := Claims{
		guid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := access.SignedString(secret)
	lenToken := len(tokenString)

	refresh := time.Now().UTC().String()
	refresh += tokenString[lenToken-5:]
	refresh = base64.StdEncoding.EncodeToString([]byte(refresh))
	bytesRefresh, err := bcrypt.GenerateFromPassword([]byte(refresh), 14)
	createRefreshToken(string(bytesRefresh), guid)
	return tokenString, refresh, err
}

func ParseAccessToken(tokenString string) (string, error) {
	access, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if claims, ok := access.Claims.(*Claims); ok && access.Valid {
		return claims.Guid, nil
	} else {
		return "", err
	}
}

func IsValidTokens(access, refreshBase64 string) (bool, string) {
	guid, err := ParseAccessToken(access)
	result := readRefreshToken(guid)
	refresh, err := base64.StdEncoding.DecodeString(refreshBase64)
	if access[len(access)-5:] == string(refresh[len(refresh)-5:]) {
		if err := bcrypt.CompareHashAndPassword([]byte(result.Refresh), []byte(refreshBase64)); err == nil {
			return true, guid
		}
	} else {
		fmt.Println("Токены не совместимы")
	}

	if err != nil {
		fmt.Println("Токен не валиден!!!")
	}
	return false, ""
}
