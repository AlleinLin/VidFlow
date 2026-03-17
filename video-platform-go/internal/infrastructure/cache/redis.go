package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/video-platform/go/internal/config"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *RedisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *RedisCache) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

func (c *RedisCache) HSet(ctx context.Context, key string, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.HSet(ctx, key, field, data).Err()
}

func (c *RedisCache) HGet(ctx context.Context, key string, field string, dest interface{}) error {
	val, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

func (c *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

func (c *RedisCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	data := make([]interface{}, len(values))
	for i, v := range values {
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data[i] = d
	}
	return c.client.LPush(ctx, key, data...).Err()
}

func (c *RedisCache) RPush(ctx context.Context, key string, values ...interface{}) error {
	data := make([]interface{}, len(values))
	for i, v := range values {
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data[i] = d
	}
	return c.client.RPush(ctx, key, data...).Err()
}

func (c *RedisCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

func (c *RedisCache) LLen(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}

func (c *RedisCache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SAdd(ctx, key, members...).Err()
}

func (c *RedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

func (c *RedisCache) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.client.SIsMember(ctx, key, member).Result()
}

func (c *RedisCache) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SRem(ctx, key, members...).Err()
}

func (c *RedisCache) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return c.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

func (c *RedisCache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

func (c *RedisCache) ZRangeByScore(ctx context.Context, key string, min, max string, offset, count int64) ([]string, error) {
	return c.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (c *RedisCache) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, key, members...).Err()
}

func (c *RedisCache) ZScore(ctx context.Context, key, member string) (float64, error) {
	return c.client.ZScore(ctx, key, member).Result()
}

func (c *RedisCache) ZRank(ctx context.Context, key, member string) (int64, error) {
	return c.client.ZRank(ctx, key, member).Result()
}

func (c *RedisCache) ZCard(ctx context.Context, key string) (int64, error) {
	return c.client.ZCard(ctx, key).Result()
}

func (c *RedisCache) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.client.Publish(ctx, channel, data).Err()
}

func (c *RedisCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}

func (c *RedisCache) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}

func (c *RedisCache) TxPipeline() redis.Pipeliner {
	return c.client.TxPipeline()
}

func (c *RedisCache) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return c.client.Watch(ctx, fn, keys...)
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) Client() *redis.Client {
	return c.client
}
