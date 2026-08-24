package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr, OrderGRPCAddr, OrderGRPCTarget, NotificationGRPCAddr, NotificationGRPCTarget, OrderDatabaseURL, NotificationDatabaseURL string
	DBMaxConns                                                                                                                        int32
}

func Load() Config {
	return Config{HTTPAddr: env("HTTP_ADDR", ":8080"), OrderGRPCAddr: env("ORDER_GRPC_ADDR", ":9090"), OrderGRPCTarget: env("ORDER_GRPC_TARGET", "localhost:9090"), NotificationGRPCAddr: env("NOTIFICATION_GRPC_ADDR", ":9091"), NotificationGRPCTarget: env("NOTIFICATION_GRPC_TARGET", "localhost:9091"), OrderDatabaseURL: env("ORDER_DATABASE_URL", "postgres://app:app@localhost:5432/orders?sslmode=disable"), NotificationDatabaseURL: env("NOTIFICATION_DATABASE_URL", "postgres://app:app@localhost:5433/notifications?sslmode=disable"), DBMaxConns: envInt("DB_MAX_CONNS", 10)}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func envInt(k string, d int32) int32 {
	v, e := strconv.ParseInt(env(k, ""), 10, 32)
	if e != nil || v < 1 {
		return d
	}
	return int32(v)
}
