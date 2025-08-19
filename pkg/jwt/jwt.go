package jwtPack

import (
	"BankApp/pkg/redisPack"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SecretService struct {
	key []byte
}

func NewSecretService(key []byte) *SecretService {
	return &SecretService{key: key}
}
func (secret *SecretService) CreateJWT(userId int32) (string, error) {
	userClaims := jwt.MapClaims{
		"sub":  fmt.Sprintf("%d", userId),
		"name": userId,
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	return userToken.SignedString(secret.key)
}

func (secret *SecretService) VerifyJWT(userToken string) (bool, error) {
	_, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return secret.key, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, fmt.Errorf("token expired")
		}
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func (secret *SecretService) decodeJWT(tokenString string) (jwt.MapClaims, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return secret.key, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func (service *SecretService) GetIdFromClaims(tokenString string) (int32, error) {
	claims, err := service.decodeJWT(tokenString)
	if err != nil {
		return -1, err
	}
	idField, exist := claims["name"]
	if !exist {
		return -1, errors.New("field name doesn't exist")
	}
	id, ok := idField.(int32)
	if !ok {
		return -1, errors.New("field name has wrong type")
	}
	return id, nil
}

func SaveJWTInRedis(ctx context.Context, token string, ID uint, time time.Duration) error {

	if err := redisPack.GetRedis().Set(ctx, token, strconv.FormatUint(uint64(ID), 10), time).Err(); err != nil {
		return err
	}
	return nil
}
