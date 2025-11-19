package redisx

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var (
	RedisIns           *redis.Client
	RedisInsOfPlatform *redis.Client
	redisConf          RedisConf
	initOnceObj        = sync.Once{}
)

// 使用该库方法前，必须先调用 Init 进行初始化
func Init(conf RedisConf) {
	if RedisIns == nil {
		redisDefaultConf := conf
		initOnceObj.Do(func() {
			RedisIns = NewClient(redisDefaultConf)
			redisConf = conf
		})
	}
}

// 获取 redis 客户端实例。使用该方法前，需确定已经调用过 Init 进行初始化
func GetRedisClient() *redis.Client {
	if RedisIns == nil {
		panic("please call Init func before using it")
	}
	_, err := RedisIns.Ping().Result()
	if err != nil || err == redis.Nil {
		RedisIns = NewClient(redisConf)
	}
	return RedisIns
}

func NewClient(conf RedisConf) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.GetAddr(),
		Password: conf.GetPwd(),
		DB:       conf.GetDb(),

		PoolSize:     1000,
		MinIdleConns: 10,
	})

	return client
}

func SetNx(key string, v any) (bool, error) {
	cli := GetRedisClient()
	isOk, err := cli.SetNX(key, v, time.Second*10).Result()
	if err != nil {
		return false, errors.Wrap(err, "setNx err")
	}
	return isOk, nil
}

func SetNxExpire(key string, v any, expr time.Duration) (bool, error) {
	cli := GetRedisClient()
	isOk, err := cli.SetNX(key, v, expr).Result()
	if err != nil {
		return false, errors.Wrap(err, "SetNxExpire err")
	}
	return isOk, nil
}

// GetLock try get lock, and expr is 10 sec
func GetLock(key string) bool {
	isOk, err := SetNx(key, 1)
	if err != nil {
		log.Printf("[GetLock] err: %v", err)
		return isOk
	}

	return isOk
}

func GetLockWithExpr(key string, expr time.Duration) bool {
	isOk, err := SetNxExpire(key, 1, expr)
	if err != nil {
		log.Printf("[GetLock] err: %v", err)
		return isOk
	}

	return isOk
}

func UnLock(key string) (bool, error) {
	cli := GetRedisClient()
	_, err := cli.Del(key).Result()
	if err != nil {
		return false, errors.Wrap(err, "unlock err")
	}
	return true, nil
}
