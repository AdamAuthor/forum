package app

import (
	"forum/internal/config"
	"forum/internal/handlers"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"forum/pkg/database"
	"log"
)

func Run() {
	config, err := config.ReadConfig("configs/config.json")
	if err != nil {
		log.Fatalf("Error with reading configs: %v", err)
	}

	db, err := database.ConnectDB(
		config.Database.Driver,
		config.Database.Path,
		config.Database.FileName,
		config.Database.SchemesDir,
	)

	if err != nil {
		log.Fatalf("Error with creating db: %v", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services)

	server := server.NewServer(":8080", handlers.InitRoutes())
	server.RunServer()
}
