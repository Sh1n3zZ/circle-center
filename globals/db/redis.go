package globals

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// RedisClient is the global Redis client
	RedisClient *redis.Client
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultRedisConfig returns default Redis configuration
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}
}

// ConnectRedis establishes a connection to Redis
func ConnectRedis(config *RedisConfig) error {
	if config == nil {
		config = DefaultRedisConfig()
	}

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	RedisClient = client
	slog.Info("Successfully connected to Redis",
		"host", config.Host,
		"port", config.Port,
		"db", config.DB,
		"pool_size", config.PoolSize,
		"min_idle_conns", config.MinIdleConns,
		"max_retries", config.MaxRetries,
		"dial_timeout", config.DialTimeout,
		"read_timeout", config.ReadTimeout,
		"write_timeout", config.WriteTimeout,
	)
	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// GetRedisClient returns the global Redis client
func GetRedisClient() *redis.Client {
	return RedisClient
}

// RedisKeyValue represents a key-value pair for Redis operations
type RedisKeyValue struct {
	Key   string
	Value interface{}
}

// RedisHash represents a hash field-value pair
type RedisHash struct {
	Field string
	Value interface{}
}

// Set sets a key-value pair in Redis
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key
func Get(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Get(ctx, key).Result()
}

// Del deletes one or more keys from Redis
func Del(ctx context.Context, keys ...string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Del(ctx, keys...).Result()
}

// Exists checks if one or more keys exist in Redis
func Exists(ctx context.Context, keys ...string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Exists(ctx, keys...).Result()
}

// Expire sets the expiration time for a key
func Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Expire(ctx, key, expiration).Result()
}

// TTL gets the remaining time to live for a key
func TTL(ctx context.Context, key string) (time.Duration, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.TTL(ctx, key).Result()
}

// Incr increments a key's value by 1
func Incr(ctx context.Context, key string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Incr(ctx, key).Result()
}

// IncrBy increments a key's value by the specified amount
func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.IncrBy(ctx, key, value).Result()
}

// HSet sets a hash field-value pair
func HSet(ctx context.Context, key string, field string, value interface{}) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.HSet(ctx, key, field, value).Err()
}

// HGet retrieves a hash field value
func HGet(ctx context.Context, key, field string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.HGet(ctx, key, field).Result()
}

// HGetAll retrieves all hash fields and values
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.HGetAll(ctx, key).Result()
}

// HDel deletes one or more hash fields
func HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.HDel(ctx, key, fields...).Result()
}

// LPush pushes one or more values to the left of a list
func LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.LPush(ctx, key, values...).Result()
}

// RPush pushes one or more values to the right of a list
func RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.RPush(ctx, key, values...).Result()
}

// LPop pops a value from the left of a list
func LPop(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.LPop(ctx, key).Result()
}

// RPop pops a value from the right of a list
func RPop(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.RPop(ctx, key).Result()
}

// LRange gets a range of elements from a list
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.LRange(ctx, key, start, stop).Result()
}

// SAdd adds one or more members to a set
func SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.SAdd(ctx, key, members...).Result()
}

// SMembers gets all members of a set
func SMembers(ctx context.Context, key string) ([]string, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.SMembers(ctx, key).Result()
}

// SRem removes one or more members from a set
func SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.SRem(ctx, key, members...).Result()
}

// ZAdd adds one or more members to a sorted set
func ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.ZAdd(ctx, key, members...).Result()
}

// ZRange gets a range of members from a sorted set
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.ZRange(ctx, key, start, stop).Result()
}

// ZRem removes one or more members from a sorted set
func ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.ZRem(ctx, key, members...).Result()
}

// Pipeline returns a new pipeline for batch operations
func Pipeline() redis.Pipeliner {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Pipeline()
}

// TxPipeline returns a new transaction pipeline
func TxPipeline() redis.Pipeliner {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.TxPipeline()
}

// Watch watches keys for modifications during a transaction
func Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Watch(ctx, fn, keys...)
}

// Multi starts a Redis transaction
func Multi() redis.Pipeliner {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.TxPipeline()
}

// FlushDB flushes the current database
func FlushDB(ctx context.Context) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.FlushDB(ctx).Err()
}

// FlushAll flushes all databases
func FlushAll(ctx context.Context) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.FlushAll(ctx).Err()
}

// Info returns Redis server information
func Info(ctx context.Context, section ...string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.Info(ctx, section...).Result()
}

// ClientList returns information about Redis clients
func ClientList(ctx context.Context) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.ClientList(ctx).Result()
}

// MemoryUsage returns memory usage of a key
func MemoryUsage(ctx context.Context, key string, samples ...int) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.MemoryUsage(ctx, key, samples...).Result()
}
