package main

import (
	"log"
	"os"
	"time"

	"github.com/xdefrag/hash-ketchum/api"
	"github.com/xdefrag/hash-ketchum/pkg/adapter/redis"
	"github.com/xdefrag/hash-ketchum/pkg/tools/config"
	"github.com/xdefrag/hash-ketchum/pkg/usecase"
)

func main() {

	// Configuration
	redisCfg := redis.Config{
		Host:                config.WithDefaultString(os.Getenv("REDIS_HOST"), "0.0.0.0"),
		Port:                config.WithDefaultInt(config.Atoi(os.Getenv("REDIS_PORT")), 6379),
		Password:            os.Getenv("REDIS_PASSWORD"),
		Database:            config.Atoi(os.Getenv("REDIS_DATABASE")),
		PoolMaxIdle:         config.Atoi(os.Getenv("REDIS_POOL_MAX_IDLE")),
		PoolMaxActive:       config.Atoi(os.Getenv("REDIS_POOL_MAX_ACTIVE")),
		PoolIdleTimeout:     time.Duration(config.Atoi(os.Getenv("REDIS_POOL_IDLE_TIMEOUT"))),
		PoolMaxConnLifetime: time.Duration(config.Atoi(os.Getenv("REDIS_POOL_MAX_CONN_LIFETIME"))),
	}

	port := config.WithDefaultInt(config.Atoi(os.Getenv("SERVER_PORT")), 8080)

	users := []string{
		config.WithDefaultString(os.Getenv("SERVER_LOGIN"), "user"),
	}

	// Dependencies
	redisLogger := log.New(os.Stdout, "REDIS ADAPTER: ", 0)

	storage := redis.New(redisCfg, redisLogger)

	usecaseHash := usecase.NewHash(storage)

	auth := NewAuth(users)

	logger := log.New(os.Stdout, "SERVER: ", 0)

	server := api.NewServer(port, auth, usecaseHash, logger)

	// Run
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

// Auth struct for map authorizer.
type Auth struct {
	users map[string]struct{}
}

// NewAuth creates Auth struct and converts users slices to map for faster access.
func NewAuth(uu []string) Auth {
	users := make(map[string]struct{})

	for _, u := range uu {
		users[u] = struct{}{}
	}

	return Auth{users}
}

// Authorize searches users map and return true if login exists.
func (a Auth) Authorize(login string) bool {
	_, ok := a.users[login]

	return ok
}
