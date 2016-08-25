package rogue

import (
	"time"

	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

func newJwt() string {

	signingKey := []byte("top secret")

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}
