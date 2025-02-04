package jwtPack

import (
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT() *jwt.Token {
	return jwt.New(jwt.SigningMethodES384)
}
