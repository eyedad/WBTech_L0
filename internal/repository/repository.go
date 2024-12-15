package repository

import (
	"example.com/m/v2/config"
	"example.com/m/v2/pkg/logging"
	postgersDB "example.com/m/v2/pkg/postgers"
	redisCache "example.com/m/v2/pkg/redis"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db    *sqlx.DB
	cache *redis.Client
}

func New(cfg *config.Config, logger *logging.Logger) *Repository {
	logger.Info("Connecting to database")
	dns := cfg.GetDNS()
	logger.Info(dns)
	db, err := postgersDB.New(dns)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to redis client")
	redis := redisCache.New(cfg.Redis.RedisHost, cfg.Redis.RedisPort, cfg.Redis.RedisDB)

	return &Repository{db, redis}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
