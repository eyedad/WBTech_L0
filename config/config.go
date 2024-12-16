package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool  `yaml:"is_bebug"`
	Logs    string `yaml:"logs"`
	Listen  struct {
		BindIP       string `yaml:"bind_ip"`
		Port         string `yaml:"port"`
		WriteTimeout int    `yaml:"write_timeout"`
		ReadTimeout  int    `yaml:"read_timeout"`
	} `yaml:"listen"`
	Postgers struct {
		DBHost   string `yaml:"db_host"`
		DBPort   string `yaml:"db_port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
	Redis struct {
		RedisHost string `yaml:"redis_host"`
		RedisPort string `yaml:"redis_port"`
		RedisDB   int    `yaml:"redis_db"`
	} `yaml:"cache"`
}

var instance *Config
var once sync.Once

func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{}
		cleanenv.ReadConfig("config/config.yml", instance)
	})
	return instance, err
}

func (cfg Config) GetDNS() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgers.DBHost,
		cfg.Postgers.DBPort,
		cfg.Postgers.Username,
		cfg.Postgers.DBName,
		cfg.Postgers.Password,
		cfg.Postgers.SSLMode,
	)
}
