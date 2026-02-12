package main

import (
	"log"
	"os"

	"github.com/danilobml/go-cached/internal/db"
	"github.com/danilobml/go-cached/internal/handlers"
	"github.com/danilobml/go-cached/internal/httpx"
	"github.com/danilobml/go-cached/internal/repositories"
	"github.com/danilobml/go-cached/internal/services"
)

func main() {
	postgresDsn := os.Getenv("POSTGRES_DSN")
	if postgresDsn == "" {
		log.Fatal("unable to read POSTGRES_DSN from env")
	}

	dbConnPool, err := db.InitDB(postgresDsn)
	if err != nil {
		log.Fatal("failed to initialize database", err)
	}
	defer dbConnPool.Close()

	userRepo := repositories.NewPgUserRepository(dbConnPool)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	httpx.StartServer(userHandler)
}
