package cache

// Cache is the basic interface for cache implementations
type Cache interface {
	Get(key string, pointer interface{}) error
	Set(key string, value interface{}) error
	Delete(key string) error
	Exists(key string) error
}
