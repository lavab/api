package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Scripts!
var (
	scriptDeleteMask = redis.NewScript(0, `
		for _, k in ipairs( redis.call( 'keys', ARGV[1] ) ) do
			redis.call( 'del', k )
		end
	`)
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

// Get retrieves data from the database and then decodes it into the passed pointer.
func (r *RedisCache) Get(key string, pointer interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()

	// Perform the get
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}

	// Initialize a new decoder
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Decode it into pointer
	return dec.Decode(pointer)
}

// Set encodes passed value and sends it to redis
func (r *RedisCache) Set(key string, value interface{}, expires time.Duration) error {
	conn := r.pool.Get()
	defer conn.Close()

	// Initialize a new encoder
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	// Encode the value
	if err := enc.Encode(value); err != nil {
		return err
	}

	// Save it into redis
	if expires == 0 {
		_, err := conn.Do("SET", key, buffer.Bytes())
		return err
	}

	_, err := conn.Do("SETEX", key, expires.Seconds(), buffer.Bytes())
	return err
}

// Delete removes data in redis by key
func (r *RedisCache) Delete(key string) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := redis.Int(conn.Do("DEL", key))
	return err
}

// DeleteMask removes data using KEYS masks
func (r *RedisCache) DeleteMask(mask string) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := scriptDeleteMask.Do(conn, mask)
	return err
}

// DeleteMulti removes multiple keys
func (r *RedisCache) DeleteMulti(keys ...interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := redis.Int(conn.Do("DEL", keys...))
	return err
}

// Exists performs a check whether a key exists
func (r *RedisCache) Exists(key string) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}
