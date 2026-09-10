package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Valkey struct {
	Client *redis.Client
}

func NewValkey() *Valkey {

	port := 6379

	if value := os.Getenv("VALKEY_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			port = parsed
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			os.Getenv("VALKEY_HOST"),
			port,
		),
		Password: os.Getenv("VALKEY_PASSWORD"),
		DB:       0,
	})

	return &Valkey{
		Client: client,
	}
}

func (v *Valkey) Ping(
	ctx context.Context,
) error {
	return v.Client.Ping(ctx).Err()
}

func (v *Valkey) BlacklistToken(
	ctx context.Context,
	jti string,
	ttl time.Duration,
) error {

	return v.Client.Set(
		ctx,
		"auth:blacklist:"+jti,
		"1",
		ttl,
	).Err()
}

func (v *Valkey) IsTokenBlacklisted(
	ctx context.Context,
	jti string,
) (bool, error) {

	result, err := v.Client.Exists(
		ctx,
		"auth:blacklist:"+jti,
	).Result()

	if err != nil {
		return false, err
	}

	return result > 0, nil
}
