package app

import (
	"context"
	"jsonjunk/config"
	"jsonjunk/internal/repository"
	"jsonjunk/internal/router"
	"jsonjunk/internal/scheduler"
	"jsonjunk/internal/service"
	logger "jsonjunk/pkg/logging"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go listenForShutdown(cancel)
	start(ctx, cfg.DBName)
}

func start(ctx context.Context, dbName string) {
	scheduler.Open(ctx)
	repo := repository.NewMongoPasteRepository(ctx, dbName)
	svc := service.NewPasteService(ctx, repo)
	router.Run(ctx, svc)
}

func listenForShutdown(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Log.Info("shutdown service...")
	cancelFunc()
}
