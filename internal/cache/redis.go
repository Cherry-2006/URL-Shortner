package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func redisAddr() string {
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		return v
	}
	return "localhost:6379"
}

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr: redisAddr(),
	})

	ctx := context.Background()

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func Get(key string) (string, error) {
	return Client.Get(context.Background(), key).Result()
}

func Set(key string, value string) {
	Client.Set(context.Background(), key, value, 0)
}

func SetWithTTL(key string, value string, ttl time.Duration) {
	Client.Set(context.Background(), key, value, ttl)
}
