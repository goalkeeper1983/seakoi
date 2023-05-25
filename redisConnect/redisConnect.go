package redisConnect

import (
	"context"
	"github.com/goalkeeper1983/seakoi/tools"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func CreateRedisConnect(redisOption ...string) *redis.Client {
	redisOptions := &redis.Options{
		Addr:     redisOption[0] + ":" + redisOption[1],
		Password: redisOption[2],
		DB:       tools.StringToInt(redisOption[3]),
	}
	redisClient := redis.NewClient(redisOptions)
	result, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", err.Error()))
	}
	tools.Log.Info(tools.RunFuncName(), zap.String(result, "Connect OK!"))
	return redisClient
}
