package route

import (
	"module/etc"
	"module/redis"
)

var redisClient *redis.RediStore

func init() {
	client, err := redis.NewRediStore(
		4,
		"tcp",
		etc.Etc.String("module/redis", "addr"),
		"", // no password set
	)
	if err != nil {
		panic(err)
	}
	redisClient = client
}
