package main

import (
	"log"
	"os"
	"strconv"

	"github.com/danilobml/go-cached/internal/cache"
	"github.com/danilobml/go-cached/internal/db"
	"github.com/danilobml/go-cached/internal/handlers"
	"github.com/danilobml/go-cached/internal/httpx"
	"github.com/danilobml/go-cached/internal/repositories"
	"github.com/danilobml/go-cached/internal/services"
)

type AppConfig struct {
	postgresDsn   string
	redisAddr     string
	redisPassword string
	redisDbStr    string
}

func main() {
	appConfig := getConfigFromEnv()

	redisDB, _ := strconv.Atoi(appConfig.redisDbStr)

	dbConnPool, err := db.InitDB(appConfig.postgresDsn)
	if err != nil {
		log.Fatal("failed to initialize database", err)
	}
	defer dbConnPool.Close()

	redisCfg := cache.Config{
		Addr:     appConfig.redisAddr,
		Password: appConfig.redisPassword,
		DB:       redisDB,
	}

	redisClient, err := cache.InitRedis(redisCfg)
	if err != nil {
		log.Fatal("failed to initialize redis", err)
	}
	defer redisClient.Close()

	userRepo := repositories.NewPgUserRepository(dbConnPool)
	userService := services.NewUserService(userRepo, redisClient)
	userHandler := handlers.NewUserHandler(userService)

	httpx.StartServer(userHandler)
}

func getConfigFromEnv() AppConfig {
	postgresDsn := os.Getenv("POSTGRES_DSN")
	if postgresDsn == "" {
		log.Fatal("unable to read POSTGRES_DSN from env")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("unable to read REDIS_ADDR from env")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		log.Fatal("unable to read REDIS_PASSWORD from env")
	}

	redisDbStr := os.Getenv("REDIS_DB")
	if redisDbStr == "" {
		log.Fatal("unable to read REDIS_DB from env")
	}

	return AppConfig{
		postgresDsn:   postgresDsn,
		redisAddr:     redisAddr,
		redisPassword: redisPassword,
		redisDbStr:    redisDbStr,
	}
}
