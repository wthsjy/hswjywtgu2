package redis

import (
	"sync"
	"strings"
	"errors"
	"github.com/go-redis/redis"
)

var redisConfig = sync.Map{}

// RegisterRedis register redis
// dsn format -> host:port:pwd, using ':' to split
func RegisterRedis(name string, dsn string) error {
	t := strings.Split(dsn, ":")

	var addr, pwd string

	switch {
	case len(t) == 3:
		addr = strings.Join(t[:2], ":")
		pwd = t[2]
	case len(t) == 2:
		addr = dsn
	case len(t) == 1:
		addr = dsn + ":6379"
	default:
		return errors.New("redis dsn format error")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       0,   // use default DB
		PoolSize: 256,
	})

	redisConfig.Store(name, client)

	return client.Ping().Err()
}

// Client CloseRedis
func Client(name string) *redis.Client {
	v, ok := redisConfig.Load(name)
	if !ok {
		return nil
	}

	c, ok := v.(*redis.Client)
	if !ok {
		return nil
	}

	return c
}

// CloseAll close all redis
func CloseAll() error {
	redisConfig.Range(func(k, v interface{}) bool {
		if c, ok := v.(*redis.Client); ok && c != nil {
			c.Close()
		}
		return false
	})
	return nil
}
