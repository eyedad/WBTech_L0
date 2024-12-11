package repository

import (
	"fmt"

	"example.com/m/v2/internal/config"
	"example.com/m/v2/pkg/logging"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config, logger *logging.Logger) *sqlx.DB {
	dns := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.DBPort, cfg.Database.Username, cfg.Database.DBName, cfg.Database.Password, cfg.Database.SSLMode)
	logger.Info(dns)
	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}
	return db
}
