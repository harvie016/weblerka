package routes

import (
	"go-web-site/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	return router
}
