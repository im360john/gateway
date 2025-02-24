package lrucache

import "time"

type Config struct {
	MaxSize int           `yaml:"max_size"`
	TTL     time.Duration `yaml:"ttl"`
}

func (c Config) Tag() string {
	return "lru_cache"
}
