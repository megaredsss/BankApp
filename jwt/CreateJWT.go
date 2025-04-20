package jwtPack

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(userId int, email string) (string, error) {
	userClaims := jwt.MapClaims{
		"sub":  fmt.Sprintf("%d", userId),
		"name": email,
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	return userToken.SignedString(getSecretKey())
}

func VerifyJWT(userToken string) (bool, error) {
	token, err := jwt.Parse(userToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong SigningMethod: %v", t.Header["alg"])
		}
		return getSecretKey(), nil
	})
	if err != nil {
		return false, fmt.Errorf("error in parsing: %v", err)
	}
	if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expTime := int64(tokenClaims["exp"].(float64))
		if time.Now().Unix() > expTime {
			return false, fmt.Errorf("token expired")
		}
		return true, nil
	}
	return false, fmt.Errorf("not valid token")
}
