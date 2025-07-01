package app

import (
	"jsonjunk/config"
	"jsonjunk/internal/repository"
	"jsonjunk/internal/router"
	"jsonjunk/internal/service"
)

func Run() {
	config.InitMongo()
	repo := repository.NewMongoPasteRepository()
	svc := service.NewPasteService(repo)
	router.Run(svc)
}
