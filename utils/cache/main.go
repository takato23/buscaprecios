package cache

import (
    "context"
    "github.com/go-redis/redis/v9"
    "ratoneando/config"
    "ratoneando/utils/logger"
)

var Client *redis.Client
var ctx = context.Background()

func Init() {
    redisURL := config.REDIS_URL
    if redisURL == "" {
        logger.LogFatal("REDIS_URL is not set, cannot initialize Redis")
    }

    opts, err := redis.ParseURL(redisURL)
    if err != nil {
        logger.LogFatal("Failed to parse REDIS_URL: " + err.Error())
    }

    Client = redis.NewClient(opts)

    // Verificar la conexión
    _, err = Client.Ping(ctx).Result()
    if err != nil {
        logger.LogFatal("Failed to connect to Redis: " + err.Error())
    }

    logger.Log("Successfully connected to Redis")
}

func Get(key string) (string, error) {
    if Client == nil {
        logger.LogWarn("Redis client not initialized, skipping cache Get")
        return "", nil
    }

    value, err := Client.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", nil
    }
    if err != nil {
        logger.LogWarn("Error getting key from Redis: " + err.Error())
        return "", err
    }
    return value, nil
}

func Set(key, value string) error {
    if Client == nil {
        logger.LogWarn("Redis client not initialized, skipping cache Set")
        return nil
    }

    err := Client.Set(ctx, key, value, 0).Err()
    if err != nil {
        logger.LogWarn("Error setting key in Redis: " + err.Error())
    }
    return err
}