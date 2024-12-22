package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TejasGhatte/fampay-task-2024/initializers"
	"github.com/redis/go-redis/v9"
)

var ctx = context.TODO()
var CacheExpirationTime = 5 * time.Minute

func GetFromCache(key string) (string, error) {
	if initializers.RedisClient == nil {
		return "", fmt.Errorf("redis client not found")
	}

	data, err := initializers.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("item not found in cache")
		}
		go LogServerError("Error Getting from cache", err, "redis")
		return "", fmt.Errorf("error getting from cache")
	}
	return data, nil
}

func SetToCache(key string, data []byte) error {
	if initializers.RedisClient == nil {
		return fmt.Errorf("redis client not found")
	}

	if err := initializers.RedisClient.Set(ctx, key, data, CacheExpirationTime).Err(); err != nil {
		go LogServerError("Error Setting to cache", err, "redis")
		return fmt.Errorf("error setting to cache")
	}
	return nil
}

func RemoveFromCache(key string) error {
	if initializers.RedisClient == nil {
		return fmt.Errorf("redis client not found")
	}

	err := initializers.RedisClient.Del(ctx, key).Err()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		go LogServerError("Error Removing from cache", err, "redis")
		return fmt.Errorf("error removing from cache")
	}
	return nil
}

func GetFromCacheGeneric(key string, model interface{}) error {
	data, err := GetFromCache(key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), model); err != nil {
		return fmt.Errorf("error while unmarshaling %s: %w", key, err)
	}

	return nil
}

func SetToCacheGeneric(key string, model interface{}) error {
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("error while marshaling %s: %w", key, err)
	}

	if err := SetToCache(key, data); err != nil {
		return err
	}

	return nil
}