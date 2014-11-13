package cache

// Setup prepares the caching environment
func Setup(redisAddr string) (*DefaultTokenCache, error) {
	redisCache, err := NewRedisCache(redisAddr)
	if err != nil {
		return nil, err
	}

	return NewTokenCache(redisCache)
}
