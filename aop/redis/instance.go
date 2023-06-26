package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"greatestworks/aop/logger"
	"time"
)

var (
	nonCacheRedis   *redis.Client
	cacheRedis      *redis.Client
	rateLimitRedis  *redis.Client
	loginQueueRedis *redis.Client
)

type Config struct {
	RankRedis               string
	RankRedisPoolSize       int
	CacheRedis              string
	CacheRedisPoolSize      int
	NonCacheRedis           string
	NonCacheRedisPoolSize   int
	RateLimitRedis          string
	RateLimitRedisPoolSize  int
	LoginQueueRedis         string
	LoginQueueRedisPoolSize int
}

type Model uint32

func InitRedisInstance(context context.Context, cfg *Config) error {
	var err error

	if cacheRedis, err = newRedisClient(context, cfg.CacheRedis, cfg.CacheRedisPoolSize); err != nil {
		return err
	}
	logger.Info("[redis] init  cache-redis client success URL:%v poolSize:%v", cfg.CacheRedis, cfg.CacheRedisPoolSize)

	if nonCacheRedis, err = newRedisClient(context, cfg.NonCacheRedis, cfg.NonCacheRedisPoolSize); err != nil {
		return err
	}
	logger.Info("[redis] init noncache-redis client success URL:%v poolSize:%v", cfg.NonCacheRedis, cfg.NonCacheRedisPoolSize)

	cfgRateLimitRedis := cfg.RateLimitRedis
	if len(cfgRateLimitRedis) == 0 {
		cfgRateLimitRedis = cfg.CacheRedis
	}
	if rateLimitRedis, err = newRedisClient(context, cfgRateLimitRedis, cfg.RateLimitRedisPoolSize); err != nil {
		return err
	}

	logger.Info("[redis] init RateLimitRedis client success URL:%v poolSize:%v", cfgRateLimitRedis, cfg.RateLimitRedisPoolSize)

	cfgLoginQueueRedis := cfg.LoginQueueRedis
	if len(cfgLoginQueueRedis) == 0 {
		cfgLoginQueueRedis = cfg.CacheRedis
	}
	if loginQueueRedis, err = newRedisClient(context, cfgLoginQueueRedis, cfg.LoginQueueRedisPoolSize); err != nil {
		return err
	}

	logger.Info("[redis] init LoginQueueRedis client success URL:%v poolSize:%v", cfgLoginQueueRedis, cfg.LoginQueueRedisPoolSize)

	return err
}

// CacheRedis ...
func CacheRedis() *redis.Client {
	return cacheRedis
}

// NonCacheRedis ...
func NonCacheRedis() *redis.Client {
	return nonCacheRedis
}

// RateLimitRedis ...
func RateLimitRedis() *redis.Client {
	return rateLimitRedis
}

// LoginQueueRedis ...
func LoginQueueRedis() *redis.Client {
	return loginQueueRedis
}

func AddDistributedLock(context context.Context, key string, val interface{}, exp time.Duration) bool {
	return nonCacheRedis.SetNX(context, key, val, exp).Val()
}

func RemDistributedLock(context context.Context, kes, args string) error {
	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	return nonCacheRedis.Eval(context, script, []string{kes}, args).Err()
}

func newRedisClient(context context.Context, url string, poolSize int) (*redis.Client, error) {

	opt, err := redis.ParseURL(url)

	if poolSize > 0 {
		opt.PoolSize = poolSize
	}

	if err != nil {
		logger.Error("new redis client url:%v error:%v", url, err)
		return nil, err
	}

	redisClient := redis.NewClient(opt)
	pong, err := redisClient.Ping(context).Result()

	if err != nil {
		logger.Error("redis client pong:%v error:%v", pong, err)
		return nil, err
	}

	logger.Debug("[redis] new redis url:%v pong:%v", url, pong)

	return redisClient, nil
}

func newRedisClusterClient(context context.Context, addrs []string) (*redis.ClusterClient, error) {

	if len(addrs) <= 0 {
		return nil, fmt.Errorf("[redis] cluster node is empty")
	}

	clusterOptions := &redis.ClusterOptions{}

	for _, addr := range addrs {
		clusterOptions.Addrs = append(clusterOptions.Addrs, addr)
	}
	redisdb := redis.NewClusterClient(clusterOptions)

	redisdb.Ping(context)

	pong, err := redisdb.Ping(context).Result()

	if err != nil {
		return nil, fmt.Errorf("redis client pong:%v error:%v", pong, err)
	}

	logger.Debug("[redis] new redis addrs:%v pong:%v", addrs, pong)

	return redisdb, nil
}
