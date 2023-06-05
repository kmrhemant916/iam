package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kmrhemant916/iam/controllers"
)


func SetupRoutes() (*chi.Mux){
	router := chi.NewRouter()
	router.Get("/", controllers.Register)

	return router
}