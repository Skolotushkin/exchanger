package config

func DefaultConfig() *Config {
	return &Config{
		KafkaBrokers:    "kafka:9092",
		KafkaTopic:      "big-transfers",
		MongoURI:        "mongodb://localhost:27017",
		MongoDBName:     "notification_db",
		MongoCollection: "transactions",
		LogLevel:        "info",
	}
}
