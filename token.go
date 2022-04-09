package main

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Claims struct {
	Guid string `json:"guid"`
	jwt.StandardClaims
}

var secret = []byte("A&'/}Z57M(2hNg=;LE?")

func GetNewRefreshToken(guid string) (string, error) {
	token := make([]byte, 10)
	for i := range token {
		token[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	token, err := bcrypt.GenerateFromPassword(token, 14)
	insertRefreshToken(string(token), guid)
	tokenStr := base64.StdEncoding.EncodeToString(token)
	return tokenStr, err
}

func GetNewAccessToken(guid string) (string, error) {
	claims := Claims{
		guid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return access.SignedString(secret)
}

func AccessTokenParse(token string) (*Claims, error) {
	access, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err == nil && access != nil {
		if claims, ok := access.Claims.(*Claims); ok && access.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func RefreshTokenValidate(access, refreshBase64 string) (bool, string) {
	claims, err := AccessTokenParse(access)
	result, err := readRefreshToken(claims.Guid)
	refresh, err := base64.StdEncoding.DecodeString(refreshBase64)
	if err := bcrypt.CompareHashAndPassword([]byte(result.Refresh), refresh); err == nil {
		return true, claims.Guid
	} else {
		fmt.Println("Токены не совпадают в БД")
	}

	if err != nil {
		fmt.Println("Токен не валиден!!!")
	}
	return false, ""
}
