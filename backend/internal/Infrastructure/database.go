package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func ConnectDB(cfg config.Config) (*bun.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi database: %w", err)
	}

	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping database: %w", err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	if cfg.AppEnv != "production" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
		))
	}

	return db, nil
}
