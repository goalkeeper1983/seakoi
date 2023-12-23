package main

import (
	"github.com/goalkeeper1983/seakoi/mysqlConnect"
	"github.com/goalkeeper1983/seakoi/redisConnect"
	"github.com/redis/go-redis/v9"
	"sync"

	"gorm.io/gorm"
)

type dbClientInstance struct {
	mysqlOnce     sync.Once
	mysqlInstance *gorm.DB
}

// GetMysqlClient 初始链接的时候需要传入参数 mysqlOption: 0user, 1pass, 2host, 3port, 4dbName, 5charset
func (This *dbClientInstance) GetMysqlClient(mysqlOption ...string) *gorm.DB {
	This.mysqlOnce.Do(func() {
		This.mysqlInstance = mysqlConnect.CreateMysqlConnect(mysqlOption...)
	})
	return This.mysqlInstance
}

var MysqlDbInstance *dbClientInstance

type redisClientInstance struct {
	redisOnce   sync.Once
	redisClient *redis.Client
}

func (This *redisClientInstance) GetRedisClient(redisOption ...string) *redis.Client {
	This.redisOnce.Do(func() {
		This.redisClient = redisConnect.CreateRedisConnect(redisOption...)
	})
	return This.redisClient
}

var RedisInstance *redisClientInstance

func main() {
	MysqlDbInstance = new(dbClientInstance)
	MysqlDbInstance.GetMysqlClient("root", "123456", "127.0.0.1", "3306", "dbname", "utf8mb4")

	RedisInstance = new(redisClientInstance)
	RedisInstance.GetRedisClient("127.0.0.1", "6379", "123456", "1")

	//log 输出
}
