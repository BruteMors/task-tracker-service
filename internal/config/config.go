package config

type Config struct {
	BindAddr    string `toml:"bind_addr" env:"bind_addr"`
	LogLevel    string `toml:"log_level" env:"log_level"`
	DatabaseURL string `toml:"database_url" env:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    "localhost:8080",
		LogLevel:    "debug",
		DatabaseURL: "host=localhost dbname=daytododatabase port=5433 user=postgres password=12345678 sslmode=disable",
	}
}
