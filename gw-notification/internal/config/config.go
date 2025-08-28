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
	KafkaBrokers    string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic      string `mapstructure:"KAFKA_TOPIC_BIG_TRANSFERS"`
	MongoURI        string `mapstructure:"MONGO_URI"`
	MongoDBName     string `mapstructure:"MONGO_DB"`
	MongoCollection string `mapstructure:"MONGO_COLLECTION"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
}

// LoadConfig загружает конфигурацию с приоритетом: файл -> env vars -> defaults
func LoadConfig() (*Config, error) {
	// Парсим флаги командной строки
	configPath := flag.String("c", "", "path to config file (optional)")
	flag.Parse()

	// Инициализируем Viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Устанавливаем значения по умолчанию
	setDefaultsFromStruct(DefaultConfig())

	// Если указан путь к конфиг-файлу, читаем из него
	if *configPath != "" {
		viper.SetConfigFile(*configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("ошибка чтения конфиг файла %s: %w", *configPath, err)
		}
		log.Printf("Загружены значения конфиг файла: %s", *configPath)
	} else {
		log.Printf("Конфиг файл не обнаружен")

		if hasEnvironmentVariables() {
			log.Printf("Используем переменные окружения для конфигурации")
		} else {
			log.Printf("Переменные окружения не найдены, используем дефолтные значения")
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось декодировать конфигурацию: %w", err)
	}

	return &cfg, nil
}

// Проверяем, заданы ли переменные окружения
func hasEnvironmentVariables() bool {
	envVars := []string{
		"KAFKA_BROKERS",
		"KAFKA_TOPIC",
		"MONGO_URI",
		"MONGO_DB",
		"MONGO_COLLECTION",
		"LOG_LEVEL",
	}
	for _, envVar := range envVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}
	return false
}

// Устанавливаем значения по умолчанию
func setDefaultsFromStruct(defaultCfg *Config) {
	viper.SetDefault("KAFKA_BROKER", defaultCfg.KafkaBrokers)
	viper.SetDefault("KAFKA_TOPIC", defaultCfg.KafkaTopic)
	viper.SetDefault("MONGO_URI", defaultCfg.MongoURI)
	viper.SetDefault("MONGO_DB", defaultCfg.MongoDBName)
	viper.SetDefault("MONGO_COLLECTION", defaultCfg.MongoCollection)
	viper.SetDefault("LOG_LEVEL", defaultCfg.LogLevel)
}
