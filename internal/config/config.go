package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr    string
	GRPCAddr    string
	GRPCTarget  string
	DatabaseURL string
	DBMaxConns  int32
}

func Load() Config {
	return Config{
		HTTPAddr:    env("HTTP_ADDR", ":8080"),
		GRPCAddr:    env("GRPC_ADDR", ":9090"),
		GRPCTarget:  env("GRPC_TARGET", "localhost:9090"),
		DatabaseURL: env("DATABASE_URL", "postgres://app:app@localhost:5432/products?sslmode=disable"),
		DBMaxConns:  envInt32("DB_MAX_CONNS", 10),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt32(key string, fallback int32) int32 {
	value, err := strconv.ParseInt(env(key, ""), 10, 32)
	if err != nil || value < 1 {
		return fallback
	}
	return int32(value)
}
