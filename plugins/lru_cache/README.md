# LRU Cache Plugin

Implements LRU (Least Recently Used) caching for query responses.

## Type
- Wrapper

## Description
Caches query responses using an LRU strategy, with configurable cache size and TTL.

## Configuration

```yaml
lru_cache:
  max_size: 1000        # Maximum number of entries in cache
  ttl: "5m"            # Time-to-live for cached entries
``` 