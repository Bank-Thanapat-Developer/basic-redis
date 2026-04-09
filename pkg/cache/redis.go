package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/config"
	"github.com/redis/go-redis/v9"
)

// ErrLocked ใช้เป็น sentinel error เพื่อให้ caller แยกแยะกรณี lock ไม่ได้
var ErrLocked = errors.New("resource is locked by another process")

type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetObject(ctx context.Context, key string, value any) error
	// Incr เพิ่มค่า counter แบบ atomic คืนค่าหลังเพิ่มแล้ว
	Incr(ctx context.Context, key string) (int64, error)
	// Expire ตั้ง TTL ให้ key ที่มีอยู่แล้ว
	Expire(ctx context.Context, key string, ttl time.Duration) error
	// AcquireLock ล็อก resource ด้วย SET NX EX คืน ErrLocked ถ้าล็อกไม่ได้
	AcquireLock(ctx context.Context, key, token string, ttl time.Duration) error
	// ReleaseLock ปล่อย lock เฉพาะเมื่อ token ตรงกัน (atomic via Lua script)
	ReleaseLock(ctx context.Context, key, token string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg config.RedisConfig) (CacheService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &redisCache{client: client}, nil
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // key ไม่มี → return empty, ไม่ใช่ error
	}
	return val, err
}

func (r *redisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *redisCache) GetObject(ctx context.Context, key string, value any) error {
	bytes, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return redis.Nil
	}
	return json.Unmarshal(bytes, value)
}

func (r *redisCache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *redisCache) AcquireLock(ctx context.Context, key, token string, ttl time.Duration) error {
	// SET key token NX EX <ttl> — atomic, ได้ lock ก็ต่อเมื่อ key ยังไม่มี
	ok, err := r.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return fmt.Errorf("acquire lock error: %w", err)
	}
	if !ok {
		return ErrLocked
	}
	return nil
}

// releaseLockScript ลบ key เฉพาะเมื่อ value ตรงกับ token ของเราเท่านั้น
// กัน process อื่นที่ได้ lock ต่อจากเราไม่ถูกลบ lock โดยเราอีกคน
var releaseLockScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`)

func (r *redisCache) ReleaseLock(ctx context.Context, key, token string) error {
	return releaseLockScript.Run(ctx, r.client, []string{key}, token).Err()
}
