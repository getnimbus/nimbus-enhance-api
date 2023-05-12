package infra

import (
	"context"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/tikivn/ultrago/u_logger"

	"nimbus-enhance-api/internal/conf"
)

type RedisClient struct {
	*goredis.Client
}

func NewRedisClient() (*RedisClient, func(), error) {
	logger := u_logger.NewLogger()
	c := goredis.NewClient(&goredis.Options{
		Addr:         conf.Config.RedisAddress,
		Password:     "",
		DB:           conf.Config.RedisDB,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
	})
	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("error creating redis client: %s", err.Error())
		return nil, nil, err
	}

	cleanup := func() {
		if err := c.Close(); err != nil {
			logger.Error(err)
		}
	}

	return &RedisClient{c}, cleanup, nil
}
