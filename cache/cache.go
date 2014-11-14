package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

// Cache is a basic key/value cache interface
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value interface{}, i int64) error
	Del(key string) error
	Exists(key string) bool
	Close() error
}

// RedisCache is the Redis implementation of Cache interface
type RedisCache struct {
	redisPool *redis.Pool
}

// NewRedisCache creates a new instance of cache with db index 0
func NewRedisCache(server string) (*RedisCache, error) {
	return NewRedisCacheWithIndex(server, 0)
}

// NewRedisCacheWithIndex Creates a new instance of cache
// with specified db index ,useful for testing
func NewRedisCacheWithIndex(server string, idx int) (*RedisCache, error) {
	pool, err := newRedisPool(server, idx)
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		redisPool: pool,
	}, nil
}

// Set fn sets the given key > value with expire
func (tc *RedisCache) Set(redisKey string, value interface{}, expiresAfter int64) error {
	conn, _ := tc.getConn()
	defer conn.Close()

	// the key should not exist already
	if tc.Exists(redisKey) {
		return fmt.Errorf("The %s token is already in cache", redisKey)
	}

	// Run set and expire functions on KEY
	err := redisPipeline(conn, func() error {
		err := conn.Send("SET", redisKey, value)
		if err != nil {
			return err
		}
		err = conn.Send("EXPIRE", redisKey, expiresAfter)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Get pulls byte stream of selected key
func (tc *RedisCache) Get(redisKey string) ([]byte, error) {
	conn, _ := tc.getConn()
	defer conn.Close()

	// the key should not exist already
	if !tc.Exists(redisKey) {
		return nil, fmt.Errorf("The %s token is not in the cache", redisKey)
	}

	// Get the key from redis
	rep, err := conn.Do("GET", redisKey)
	if err != nil {
		return nil, err
	}

	byteVal, err := redis.Bytes(rep, nil)
	if err != nil {
		return nil, err
	}

	return byteVal, nil

}

// Del removes the key from Redis
func (tc *RedisCache) Del(redisKey string) error {
	conn, _ := tc.getConn()
	defer conn.Close()

	// the key should not exist already
	if !tc.Exists(redisKey) {
		return fmt.Errorf("The %s token is not in the cache", redisKey)
	}

	//now delete it from redis
	_, err := conn.Do("DEL", redisKey)
	if err != nil {
		return err
	}

	return nil

}

// Exists check for presence of given key
func (tc *RedisCache) Exists(key string) bool {
	conn, _ := tc.getConn()
	defer conn.Close()

	reply, err := conn.Do("EXISTS", key)
	if err != nil {
		return false
	}

	replyInt, err := redis.Int(reply, nil)
	return replyInt == 1
}

// Close releases the resources
func (tc *RedisCache) Close() error {
	tc.redisPool.Close()
	return nil
}

// the private methods

//We need this step to be able to inject select db
func (tc *RedisCache) getConn() (redis.Conn, error) {
	conn := tc.redisPool.Get()
	return conn, nil
}

//some private utils

// Creates a new pool
func newRedisPool(server string, idx int) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			//select the db
			_, err = c.Do("SELECT", idx)
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}

// simple utility that allows executing commands at once aka pipeline
func redisPipeline(conn redis.Conn, pipeFn func() error) error {
	//create the multi part
	_ = conn.Send("MULTI")
	err := pipeFn()
	if err != nil {
		return err
	}

	//finally you do exec
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
