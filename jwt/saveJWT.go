package jwtPack

import (
	"BankApp/pkg/redisPack"
	"context"
	"strconv"
	"time"
)

func SaveJWTInRedis(ctx context.Context, token string, ID uint, time time.Duration) error {
	if err := redisPack.GetRedis().Set(ctx, strconv.FormatUint(uint64(ID), 10), token, time).Err(); err != nil {
		return err
	}
	return nil
}
