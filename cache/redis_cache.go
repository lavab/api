package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// RedisCache is an implementation of Cache that uses Redis as a backend
type RedisCache struct {
	pool *redis.Pool
}

// RedisCacheOpts is used to pass options to NewRedisCache
type RedisCacheOpts struct {
	Address     string
	Database    int
	Password    string
	MaxIdle     int
	IdleTimeout time.Duration
}

// NewRedisCache creates a new cache with a redis backend
func NewRedisCache(options *RedisCacheOpts) (*RedisCache, error) {
	// Default values
	if options.MaxIdle == 0 {
		options.MaxIdle = 3
	}
	if options.IdleTimeout == 0 {
		options.IdleTimeout = 240 * time.Second
	}

	// Create a new redis pool
	pool := &redis.Pool{
		MaxIdle:     options.MaxIdle,
		IdleTimeout: options.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", options.Address)
			if err != nil {
				return nil, err
			}

			if options.Password != "" {
				if _, err := c.Do("AUTH", options.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if options.Database != 0 {
				if _, err := c.Do("SELECT", options.Database); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// Test the pool
	conn := pool.Get()
	defer conn.Close()
	if err := pool.TestOnBorrow(conn, time.Now()); err != nil {
		return nil, err
	}

	// Return a new cache struct
	return &RedisCache{
		pool: pool,
	}, nil
}
