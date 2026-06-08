package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	if cfg.RedisHost == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       0, // database default
	})

	// Pengecekan koneksi awal dengan batas waktu 3 detik
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal koneksi ke Redis di %s: %w", addr, err)
	}

	return rdb, nil
}
