package main

import (
	"sync"

	"github.com/goalkeeper1983/seakoi/mysqlConnect"
	"github.com/goalkeeper1983/seakoi/postgreSQLConnect"
	"github.com/goalkeeper1983/seakoi/redisConnect"
	"github.com/goalkeeper1983/seakoi/tools"
	"github.com/redis/go-redis/v9"

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

type pgsqlClientInstance struct {
	pgsqlOnce     sync.Once
	pgsqlInstance *gorm.DB
}

// host user password dbname port
func (This *pgsqlClientInstance) GetPgsqlClient(option ...string) *gorm.DB {
	This.pgsqlOnce.Do(func() {
		This.pgsqlInstance = postgreSQLConnect.CreatePostgreSQLConnect(option...)
	})
	return This.pgsqlInstance
}

func main() {
	tools.Log.Info("test")
}
