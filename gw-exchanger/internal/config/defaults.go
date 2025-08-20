package config

func DefaultConfig() *Config {
	return &Config{
		GRPCHost: "0.0.0.0",
		GRPCPort: "50051",

		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "exchanger",
		DBPassword: "exchanger_pass",
		DBName:     "exchanger_db",
		DBSSL:      "disable",

		LogLevel: "info",
	}
}
