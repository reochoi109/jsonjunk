package app

import (
	"jsonjunk/config"
	"jsonjunk/internal/repository"
	"jsonjunk/internal/router"
	"jsonjunk/internal/service"
	logger "jsonjunk/pkg/logging"
)

func setup() {
	logger.Init(false)
	config.InitMongo()
}

func Run() {
	setup()

	repo := repository.NewMongoPasteRepository()
	svc := service.NewPasteService(repo)
	router.Run(svc)
}
