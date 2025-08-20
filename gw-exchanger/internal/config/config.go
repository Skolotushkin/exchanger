package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCHost string `mapstructure:"GRPC_HOST"`
	GRPCPort string `mapstructure:"GRPC_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSL      string `mapstructure:"DB_SSL"`

	LogLevel string `mapstructure:"LOG_LEVEL"`
}

// LoadConfig загружает конфигурацию с приоритетом: файл -> env vars -> defaults
func LoadConfig() (*Config, error) {
	// Парсим флаги командной строки
	configPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	// Инициализируем Viper
	viper.AutomaticEnv() // Всегда читаем из переменных окружения
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Устанавливаем значения по умолчанию
	setDefaultsFromStruct(DefaultConfig())

	// Если указан путь к конфиг-файлу, пытаемся загрузить из него
	if *configPath != "" {
		viper.SetConfigFile(*configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("ошибка чтения конфиг файла %s: %w", *configPath, err)
		}
		log.Printf("Загружены значения конфиг файла: %s", *configPath)
	} else {
		// Если файл не указан, используем только env vars и defaults
		log.Printf("Конфиг файл не обнаружен")

		// Проверяем переменные окружения
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

// hasEnvironmentVariables проверяет, установлены ли какие-либо переменные окружения
func hasEnvironmentVariables() bool {
	envVars := []string{
		"GRPC_HOST", "GRPC_PORT",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL",
		"LOG_LEVEL",
	}

	for _, envVar := range envVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}
	return false
}

// setDefaultsFromStruct устанавливает значения по умолчанию из структуры
func setDefaultsFromStruct(defaultCfg *Config) {
	viper.SetDefault("GRPC_HOST", defaultCfg.GRPCHost)
	viper.SetDefault("GRPC_PORT", defaultCfg.GRPCPort)
	viper.SetDefault("DB_HOST", defaultCfg.DBHost)
	viper.SetDefault("DB_PORT", defaultCfg.DBPort)
	viper.SetDefault("DB_USER", defaultCfg.DBUser)
	viper.SetDefault("DB_PASSWORD", defaultCfg.DBPassword)
	viper.SetDefault("DB_NAME", defaultCfg.DBName)
	viper.SetDefault("DB_SSL", defaultCfg.DBSSL)
	viper.SetDefault("LOG_LEVEL", defaultCfg.LogLevel)
}
