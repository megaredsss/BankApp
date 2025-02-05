package jwtPack

import (
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT() *jwt.Token {
	return jwt.New(jwt.SigningMethodES384)
}

func VerifyJWT(token *jwt.Token) {
	if !token.Valid {
		return
	}
}
