package redis

import (
	"api/config"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)


var ctx = context.Background()
var RedisClient  *redis.Client

type RedisService interface {
	StateSub() 
}


func NewRedisClient() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.C.Redis.HOST, config.C.Redis.PORT)})
	
	StateSub()
	
	return RedisClient
}



