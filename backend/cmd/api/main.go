package main

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	infrastructure "github.com/ahmadzakyarifin/schoolpay/internal/Infrastructure"
	"github.com/ahmadzakyarifin/schoolpay/internal/app"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/rs/zerolog/log"

	_ "github.com/ahmadzakyarifin/schoolpay/docs"
)

// @title SchoolPay API
// @version 1.0
// @description Backend API for SchoolPay
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email ahmadzakyarifin@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()
	utils.InitLogger(cfg.AppEnv)

	db, err := infrastructure.ConnectDB(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	app := app.NewApp(db, cfg)

	port := ":" + app.Cfg.Port
	log.Info().Msgf("server running on %s", port)

	if err := app.Server.Run(port); err != nil {
		log.Fatal().Err(err).Msg("Server crashed")
	}
}
