package repository

import (
	"example.com/m/v2/config"
	"example.com/m/v2/pkg/kafka"
	"example.com/m/v2/pkg/logging"
	postgersDB "example.com/m/v2/pkg/postgres"
	redisCache "example.com/m/v2/pkg/redis"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db       *sqlx.DB
	cache    *redis.Client
	producer *kafka.Producer
}

func New(cfg *config.Config, logger *logging.Logger) *Repository {
	logger.Info("Connecting to database")
	dns := cfg.GetDNS()
	logger.Info(dns)
	db, err := postgersDB.New(dns)
	if err != nil {
		logger.Fatalf("Faild to connect to ddatabase. error: %v", err)
	}

	logger.Info("Connecting to redis client")
	redis := redisCache.New(cfg.Redis.RedisHost, cfg.Redis.RedisPort, cfg.Redis.RedisDB)

	logger.Info("Connecting to kafka")
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		logger.Errorf("Faild to connect to kafka. error: %v", err)
	}

	return &Repository{db, redis, producer}
}

func (r *Repository) Close() error {
	r.producer.Close()

	err := r.cache.Close()
	if err != nil {
		return err
	}
	err = r.db.Close()
	if err != nil {
		return err
	}

	return err
}
