package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xdefrag/hash-ketchum/pkg/tools/config"
	"github.com/xdefrag/hash-ketchum/pkg/types"
)

// Config for redis client connection.
type Config struct {
	Host     string
	Port     int
	Password string
	Database int

	PoolMaxIdle         int
	PoolMaxActive       int
	PoolIdleTimeout     time.Duration
	PoolMaxConnLifetime time.Duration
}

// Redis struct.
type Redis struct {
	cfg    Config
	pool   *redis.Pool
	logger *log.Logger
}

// New dial redis and returns Redis struct.
func New(cfg Config, logger *log.Logger) *Redis {
	cfg.Host = config.WithDefaultString(cfg.Host, "0.0.0.0")
	cfg.Port = config.WithDefaultInt(cfg.Port, 6379)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var opts []redis.DialOption

	if cfg.Password != "" {
		opts = append(opts, redis.DialPassword(cfg.Password))
	}

	if cfg.Database != 0 {
		opts = append(opts, redis.DialDatabase(cfg.Database))
	}

	pool := &redis.Pool{
		Dial:            connect(addr, opts),
		MaxIdle:         cfg.PoolMaxIdle,
		MaxActive:       cfg.PoolMaxActive,
		IdleTimeout:     cfg.PoolIdleTimeout,
		MaxConnLifetime: cfg.PoolMaxConnLifetime,
		Wait:            false,
	}

	return &Redis{cfg, pool, logger}
}

// Store saves hash in redis with format like
// map[login][]map[hash]timestamp
func (r Redis) Store(ctx context.Context, hash types.Hash) error {
	conn, err := r.pool.GetContext(ctx)
	defer conn.Close()
	if err != nil {
		return err
	}

	r.log(fmt.Sprintf("Saving new hash: %s - %s - %d", hash.Login, hash.Hash, hash.Timestamp))

	result, err := conn.Do("HMSET", hash.Login, hash.Hash, hash.Timestamp)
	if err != nil {
		r.log(err.Error())

		return err
	}

	if result.(string) != "OK" {
		r.log(err.Error())

		return fmt.Errorf("Redis error: %s", result)
	}

	return nil
}

func (r Redis) log(log string) {
	if r.logger != nil {
		r.logger.Println(log)
	}
}

func connect(addr string, opts []redis.DialOption) func() (redis.Conn, error) {
	return func() (redis.Conn, error) { return redis.Dial("tcp", addr, opts...) }
}
