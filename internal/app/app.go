package app

import (
	"jsonjunk/config"
	"jsonjunk/internal/repository"
	"jsonjunk/internal/router"
	"jsonjunk/internal/service"
	logger "jsonjunk/pkg/logging"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Run() {
	_ = godotenv.Load()
	cfg := config.LoadEnvConfig()
	logger.Init(cfg)

	logger.Log.Info("Starting service...",
		zap.String("mode", cfg.ServiceMode),
		zap.String("port", cfg.Port),
	)

	if err := config.InitMongo(cfg); err != nil {
		logger.Log.Fatal("Mongo Init Failed", zap.Error(err))
	}

	logger.Log.Debug("Initializing db...",
		zap.String("host", cfg.DBHost),
		zap.String("port", cfg.DBPort),
		zap.String("mongo_url", cfg.MongoURI),
	)

	start(cfg.DBName)
}

func start(dbName string) {
	repo := repository.NewMongoPasteRepository(dbName)
	svc := service.NewPasteService(repo)
	router.Run(svc)
}
