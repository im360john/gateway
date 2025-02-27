package lrucache

import "time"

// Config represents LRU cache configuration
type Config struct {
	// MaxSize is the maximum number of entries the cache can hold
	MaxSize int `yaml:"max_size"`

	// TTL specifies how long entries should remain in cache
	// Format: time.Duration string ("5m", "1h", "24h")
	TTL time.Duration `yaml:"ttl"`
}

func (c Config) Tag() string {
	return "lru_cache"
}

func (c Config) Doc() string {
	return docString
}
