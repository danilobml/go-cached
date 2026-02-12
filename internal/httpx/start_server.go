package httpx

import (
	"log"
	"net/http"

	"github.com/danilobml/go-cached/internal/handlers"
	"github.com/danilobml/go-cached/internal/routes"
)

const httpPort = ":8080"

func StartServer(handler *handlers.UserHandler) {
	routeHandler := routes.RegisterRoutes(handler)

	srv := http.Server{
		Addr:    httpPort,
		Handler: routeHandler,
	}

	log.Printf("Server listening on port: %s", httpPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to initialize server", err)
	}
}
