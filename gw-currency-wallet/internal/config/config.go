package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host           string        `mapstructure:"HTTP_HOST"`
		Port           string        `mapstructure:"HTTP_PORT"`
		ReadTimeout    time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
		WriteTimeout   time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
		IdleTimeout    time.Duration `mapstructure:"HTTP_IDLE_TIMEOUT"`
		MaxHeaderBytes int           `mapstructure:"HTTP_MAX_HEADER_BYTES"`
	} `mapstructure:",squash"`

	DB struct {
		Host            string        `mapstructure:"DB_HOST"`
		Port            string        `mapstructure:"DB_PORT"`
		User            string        `mapstructure:"DB_USER"`
		Password        string        `mapstructure:"DB_PASSWORD"`
		Name            string        `mapstructure:"DB_NAME"`
		SSL             string        `mapstructure:"DB_SSL"`
		MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
		MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
		ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
	} `mapstructure:",squash"`

	Redis struct {
		RedisEnabled bool   `mapstructure:"REDIS_ENABLED"`
		Addr         string `mapstructure:"REDIS_ADDR"`
		// Password     string `mapstructure:"REDIS_PASSWORD"`
		DB int `mapstructure:"REDIS_DB"`
	} `mapstructure:",squash"`

	ExchangeService struct {
		Addr string `mapstructure:"EXCHANGER_GRPC_ADDR"`
	} `mapstructure:",squash"`

	Kafka struct {
		Brokers string `mapstructure:"KAFKA_BROKERS"`
		Topic   string `mapstructure:"KAFKA_TOPIC_BIG_TRANSFERS"`
	} `mapstructure:",squash"`

	Auth struct {
		JWTSecret string `mapstructure:"JWT_SECRET"`
		JWTTTL    int    `mapstructure:"JWT_TTL"`
	} `mapstructure:",squash"`

	LogLevel string `mapstructure:"LOG_LEVEL"`
}

// LoadConfig загружает конфигурацию
func LoadConfig() (*Config, error) {
	configPath := flag.String("c", "", "path to config file (optional)")
	flag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaultsFromStruct(DefaultConfig())

	if *configPath != "" {
		viper.SetConfigFile(*configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("ошибка чтения конфиг файла %s: %w", *configPath, err)
		}
		log.Printf("Загружены значения конфиг файла: %s", *configPath)
	} else {
		log.Printf("Конфиг файл не обнаружен")
		if hasEnvironmentVariables() {
			log.Printf("Using environment variables for configuration")
		} else {
			log.Printf("No environment variables found. Using default configuration")
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}

func hasEnvironmentVariables() bool {
	envVars := []string{
		"HTTP_HOST", "HTTP_PORT", "HTTP_READ_TIMEOUT", "HTTP_WRITE_TIMEOUT", "HTTP_IDLE_TIMEOUT", "HTTP_MAX_HEADER_BYTES",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME",
		"REDIS_ENABLED", "REDIS_ADDR", "REDIS_DB",
		"EXCHANGER_GRPC_ADDR",
		"KAFKA_BROKERS", "KAFKA_TOPIC_BIG_TRANSFERS",
		"JWT_SECRET", "JWT_TTL",
		"LOG_LEVEL",
	}
	for _, envVar := range envVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}
	return false
}

func setDefaultsFromStruct(defaultCfg *Config) {
	// Server defaults
	viper.SetDefault("HTTP_HOST", defaultCfg.Server.Host)
	viper.SetDefault("HTTP_PORT", defaultCfg.Server.Port)
	viper.SetDefault("HTTP_READ_TIMEOUT", defaultCfg.Server.ReadTimeout)
	viper.SetDefault("HTTP_WRITE_TIMEOUT", defaultCfg.Server.WriteTimeout)
	viper.SetDefault("HTTP_IDLE_TIMEOUT", defaultCfg.Server.IdleTimeout)
	viper.SetDefault("HTTP_MAX_HEADER_BYTES", defaultCfg.Server.MaxHeaderBytes)

	// DB defaults
	viper.SetDefault("DB_HOST", defaultCfg.DB.Host)
	viper.SetDefault("DB_PORT", defaultCfg.DB.Port)
	viper.SetDefault("DB_USER", defaultCfg.DB.User)
	viper.SetDefault("DB_PASSWORD", defaultCfg.DB.Password)
	viper.SetDefault("DB_NAME", defaultCfg.DB.Name)
	viper.SetDefault("DB_SSL", defaultCfg.DB.SSL)
	viper.SetDefault("DB_MAX_OPEN_CONNS", defaultCfg.DB.MaxOpenConns)
	viper.SetDefault("DB_MAX_IDLE_CONNS", defaultCfg.DB.MaxIdleConns)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", defaultCfg.DB.ConnMaxLifetime)

	// Redis defaults
	viper.SetDefault("REDIS_ENABLED", defaultCfg.Redis.RedisEnabled)
	viper.SetDefault("REDIS_ADDR", defaultCfg.Redis.Addr)
	// viper.SetDefault("REDIS_PASSWORD", defaultCfg.Redis.Password)
	viper.SetDefault("REDIS_DB", defaultCfg.Redis.DB)

	// Exchange service defaults
	viper.SetDefault("EXCHANGER_GRPC_ADDR", defaultCfg.ExchangeService.Addr)

	// Kafka defaults
	viper.SetDefault("KAFKA_BROKERS", defaultCfg.Kafka.Brokers)
	viper.SetDefault("KAFKA_TOPIC_BIG_TRANSFERS", defaultCfg.Kafka.Topic)

	// Auth defaults
	viper.SetDefault("JWT_SECRET", defaultCfg.Auth.JWTSecret)
	viper.SetDefault("JWT_TTL", defaultCfg.Auth.JWTTTL)

	// Log level
	viper.SetDefault("LOG_LEVEL", defaultCfg.LogLevel)
}
