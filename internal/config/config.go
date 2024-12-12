package config

import (
	"sync"

	"example.com/m/v2/pkg/logging"
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
	Database struct {
		Host     string `yaml:"host"`
		DBPort   string `yaml:"db_port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
