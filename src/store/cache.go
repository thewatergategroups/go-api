package store

import (
	"context"
	"errors"
	cfg "go-api/src/config"
	"sync"
	"time"
)


type CacheType int
const (
	MemoryCacheType CacheType = iota
	RedisCacheType
)
var (
	CacheTypes = map[string]CacheType{
		"memory": MemoryCacheType,
		"redis": RedisCacheType,
	}
	once sync.Once
	cache CacheInterface
)


type CacheInterface interface {
	Get(key string) (string, error)
	Set(key string, value string,expiration time.Duration) error
	Del(key string) error
	HSet(headKey string, key string, value string,expiration time.Duration) error
	HGet(headKey string, key string) (string, error)
	HDel(headKey string, key string) error
	HGetAll(headKey string) (map[string]string, error)
}


func Cache() CacheInterface {
	once.Do(func(){
		cache =  NewMemoryCache()
		c, ok :=CacheTypes[cfg.Cfg().Cache]
		if !ok{
			return
		}
		if c == RedisCacheType{
			cfg.Red()
			cache = NewRedisCache()
		} 
	})
	return cache
}


type RedisCache struct {}

func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func (red *RedisCache) Get(key string) (string, error) {
	return cfg.Red().Get(context.Background(),key).Result()
}

func (red *RedisCache) Set(key string, value string,expiration time.Duration) (error) {
	_, err :=  cfg.Red().Set(context.Background(), key, value, expiration).Result()
	return err
}

func (red *RedisCache) Del(key string) error {
	_,err := cfg.Red().Del(context.Background(),key).Result()
	return err
}

func (red *RedisCache) HGet(headKey string, key string) (string, error) {
	return cfg.Red().HGet(context.Background(),headKey, key).Result()
}

func (red *RedisCache) HSet(headKey string, key string, value string,expiration time.Duration) (error) {
	_, err :=  cfg.Red().HSet(context.Background(), headKey, key, value, expiration).Result()
	return err
}

func (red *RedisCache) HDel(headKey string, key string) error {
	_,err := cfg.Red().HDel(context.Background(),headKey, key).Result()
	return err
}

func (red *RedisCache) HGetAll(headKey string) (map[string]string, error) {
	return cfg.Red().HGetAll(context.Background(),headKey).Result()
}



type MemoryCache struct {
	mu   sync.RWMutex
	data map[string]string
	hash map[string]map[string]string

}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]string),
		hash: make(map[string]map[string]string),
	}
}

func (m *MemoryCache) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return "", errors.New("failed to fetch item")
	}
	return val, nil
}

func (m *MemoryCache) Set(key string, value string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	// Note: expiration is ignored here for simplicity
	return nil
}

func (m *MemoryCache) Del(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryCache) HSet(headKey string, key string, value string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.hash[headKey]; !exists {
		m.hash[headKey] = make(map[string]string)
	}
	m.hash[headKey][key] = value
	// Ignoring expiration for now
	return nil
}

func (m *MemoryCache) HGet(headKey string, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if fields, ok := m.hash[headKey]; ok {
		if val, found := fields[key]; found {
			return val, nil
		}
	}
	return "", errors.New("key not found in hash")
}

func (m *MemoryCache) HDel(headKey string, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fields, ok := m.hash[headKey]; ok {
		delete(fields, key)
	}
	return nil
}

func (m *MemoryCache) HGetAll(headKey string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fields, ok := m.hash[headKey]
	if !ok {
		return nil, errors.New("hash not found")
	}
	// Return a copy to avoid race conditions
	result := make(map[string]string)
	for k, v := range fields {
		result[k] = v
	}
	return result, nil
}