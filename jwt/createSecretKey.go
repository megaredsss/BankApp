package jwtPack

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

var (
	secretKey []byte
	once      sync.Once
)

func InitSecretKey() {
	once.Do(func() {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			panic("Error during generating secret key")
		}
		secretKey = key
	})
	base64.StdEncoding.EncodeToString(secretKey)
}

func getSecretKey() []byte {
	InitSecretKey()
	return secretKey
}
