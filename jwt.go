package rogue

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func newJwt() string {

	token := jwt.New(jwt.SigningMethodHS256)

	signingKey := []byte("top secret")

	//token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}
