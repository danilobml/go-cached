package routes

import (
	"net/http"

	"github.com/danilobml/go-cached/internal/handlers"
)

func RegisterRoutes(handler *handlers.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", handler.GetUser)
	mux.HandleFunc("GET /users", handler.GetAllUsers)
	mux.HandleFunc("POST /users", handler.CreateNewUser)

	return mux
}
