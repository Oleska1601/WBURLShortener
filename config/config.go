package config

import (
	"time"

	"github.com/wb-go/wbf/config"
)

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type PostgresConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DB              string        `mapstructure:"db"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	Database int           `mapstructure:"database"`
	TTL      time.Duration `mapstructure:"ttl"`
}

type Config struct {
	App    AppConfig      `mapstructure:"app"`
	Server ServerConfig   `mapstructure:"server"`
	DB     PostgresConfig `mapstructure:"postgres"`
	Logger LoggerConfig   `mapstructure:"logger"`
	Redis  RedisConfig    `mapstructure:"redis"`
}

func New() (*Config, error) {
	cfg := config.New()
	cfg.LoadConfigFiles("./config/config.yaml")

	// Включить env переменные с приставкой
	cfg.EnableEnv("")

	var config Config
	err := cfg.Unmarshal(&config)
	return &config, err
}
