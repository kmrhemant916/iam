package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kmrhemant916/iam/controllers"
	"github.com/kmrhemant916/iam/database"
)


func SetupRoutes() (*chi.Mux){
	db,err := database.Connection()
	if err != nil {
		panic(err)
	}
	app := &controllers.App{
		DB: db,
	}
	router := chi.NewRouter()
	router.Post("/signup", app.Signup)

	return router
}