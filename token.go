package main

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.StandardClaims
	uuid string `json:"uuid"`
}

func createJWT() (string, error) {
	claims := Claims{}
	claims.uuid = uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.At(time.Now()),
		},
		uuid: uuid.New().String(),
	})
	return token.SignedString("secret_key")
}

/*func parseJWT(accessToken string, signingKey []byte) (string, error) {

	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHS512); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.uuid, nil
	}

	return "", jwt.ErrSignatureInvalid
}

func refreshJWT() {

}*/
