package config

import (
	"time"
)

func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Host           string        `mapstructure:"HTTP_HOST"`
			Port           string        `mapstructure:"HTTP_PORT"`
			ReadTimeout    time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
			WriteTimeout   time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
			IdleTimeout    time.Duration `mapstructure:"HTTP_IDLE_TIMEOUT"`
			MaxHeaderBytes int           `mapstructure:"HTTP_MAX_HEADER_BYTES"`
		}{
			Host:           "localhost",
			Port:           "8080",
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			IdleTimeout:    15 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1MB
		},

		DB: struct {
			Host            string        `mapstructure:"DB_HOST"`
			Port            string        `mapstructure:"DB_PORT"`
			User            string        `mapstructure:"DB_USER"`
			Password        string        `mapstructure:"DB_PASSWORD"`
			Name            string        `mapstructure:"DB_NAME"`
			SSL             string        `mapstructure:"DB_SSL"`
			MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
			MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
			ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
		}{
			Host:            "localhost",
			Port:            "5432",
			User:            "postgres",
			Password:        "postgres",
			Name:            "currency_wallet",
			SSL:             "disable",
			MaxOpenConns:    100,
			MaxIdleConns:    20,
			ConnMaxLifetime: 5 * time.Minute,
		},

		Redis: struct {
			RedisEnabled bool   `mapstructure:"REDIS_ENABLED"`
			Addr         string `mapstructure:"REDIS_ADDR"`
			// Password     string `mapstructure:"REDIS_PASSWORD"`
			DB int `mapstructure:"REDIS_DB"`
		}{
			RedisEnabled: false,
			Addr:         "localhost:6379",
			// Password:     "",
			DB: 0,
		},

		ExchangeService: struct {
			Addr string `mapstructure:"EXCHANGER_GRPC_ADDR"`
		}{
			Addr: "localhost:50051",
		},

		Kafka: struct {
			Brokers string `mapstructure:"KAFKA_BROKERS"`
			Topic   string `mapstructure:"KAFKA_TOPIC_BIG_TRANSFERS"`
		}{
			Brokers: "localhost:9092",
			Topic:   "big_transfers",
		},

		Auth: struct {
			JWTSecret string `mapstructure:"JWT_SECRET"`
			JWTTTL    int    `mapstructure:"JWT_TTL"`
		}{
			JWTSecret: "secret",
			JWTTTL:    3600,
		},

		LogLevel: "info",
	}
}
