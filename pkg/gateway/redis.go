package gateway

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

func NewRedisCilent() *redis.Client {
	env := os.Getenv("APP_ENV")
	client := &redis.Client{}
	switch env {
	case "pro":
		REDIS_URL := os.Getenv("REDIS_URL")
		option, _ := redis.ParseURL(REDIS_URL)
		client = redis.NewClient(option)
	case "dev":
		client = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})
	}
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("failed to connect to %s redis: %v", env, err)
	}
	log.Printf("%s redis client connected", env)
	return client
}
