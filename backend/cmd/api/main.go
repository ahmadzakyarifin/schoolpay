package main

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	infrastructure "github.com/ahmadzakyarifin/schoolpay/internal/Infrastructure"
	"github.com/ahmadzakyarifin/schoolpay/internal/app"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig()
	utils.InitLogger(cfg.AppEnv)

	db, err := infrastructure.ConnectDB(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	redisClient, err := infrastructure.ConnectRedis(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	log.Info().Msg("Connected to Redis successfully")
	defer redisClient.Close()

	app := app.NewApp(db, cfg, redisClient)

	port := ":" + app.Cfg.Port
	log.Info().Msgf("server running on %s", port)

	if err := app.Server.Run(port); err != nil {
		log.Fatal().Err(err).Msg("Server crashed")
	}
}
