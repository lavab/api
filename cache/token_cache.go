package cache

import (
	"encoding/json"
	"fmt"
	"github.com/lavab/api/models"
	"time"
)

const (
	tokenKey = "token:%s"
)

// TokenCache is the interface for caching the tokens in system
type TokenCache interface {
	//Gets the token from store
	GetToken(key string) (*models.Token, error)
	//Saves the token into store
	SetToken(*models.Token) error
	//Removes the token from cache
	InvalidateToken(key string) error
}

// DefaultTokenCache is the redis implementation of TokenCache
type DefaultTokenCache struct {
	Cache
}

// NewTokenCache creates a new instance of cache with db index 0
func NewTokenCache(cache Cache) (*DefaultTokenCache, error) {
	return &DefaultTokenCache{
		Cache: cache,
	}, nil

}

// SetToken sets the given model into store
func (tc *DefaultTokenCache) SetToken(token *models.Token) error {
	// generate the key
	redisKey := fmt.Sprintf(tokenKey, token.ID)

	// get the left time
	now := time.Now().UTC()
	expiresAfter := token.Expiring.ExpiryDate.Sub(now).Seconds()

	//Marshall it as json
	tokenAsBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	//Call the underlying interface
	if err := tc.Set(redisKey, tokenAsBytes, int64(expiresAfter)); err != nil {
		return err
	}
	return nil
}

// GetToken gets a token from db
func (tc *DefaultTokenCache) GetToken(key string) (*models.Token, error) {
	tokenBytes, err := tc.Get(key)
	if err != nil {
		return nil, err
	}

	//unmarshall the value here
	token := new(models.Token)
	err = json.Unmarshal(tokenBytes, token)
	if err != nil {
		return nil, fmt.Errorf("Unmarshall error : %s when pulling from cache", key)
	}

	return token, nil

}

// InvalidateToken removes the key from Redis
func (tc *DefaultTokenCache) InvalidateToken(key string) error {
	if err := tc.Delete(key); err != nil {
		return err
	}
	return nil
}
