package main

import (
	"context"
	"fmt"

	"github.com/Akrom0181/Auth/api"
	"github.com/Akrom0181/Auth/config"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/Akrom0181/Auth/storage/postgres"
	_ "github.com/joho/godotenv"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	// Postgres storage init
	store, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		fmt.Println("error while connecting to db:", err)
		return
	}
	defer store.CloseDB()

	// API server init (only using storage, no service layer)
	server := api.New(&store, log, cfg)

	fmt.Println("Program is running on localhost:8080...")
	if err := server.Run(":8080"); err != nil {
		log.Error("Failed to run server", logger.Error(err))
	}
}
