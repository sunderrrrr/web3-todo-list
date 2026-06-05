package config

import "os"

type Config struct {
	RPC      string
	DB       string
	Token    string
	TodoList string
	Port     string
}

func FromEnv() Config {
	return Config{
		RPC:      getEnv("RPC_URL", "ws://localhost:8545"),
		DB:       getEnv("DB_DSN", "postgres://indexer:indexer@localhost:5432/indexer?sslmode=disable"),
		Token:    getEnv("TOKEN_ADDR", "0x5FbDB2315678afecb367f032d93F642f64180aa3"),
		TodoList: getEnv("TODO_ADDR", "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"),
		Port:     getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
