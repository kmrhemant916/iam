package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kmrhemant916/iam/controllers"
	"gorm.io/gorm"
)


func SetupRoutes(db *gorm.DB) (*chi.Mux){
	app := &controllers.App{
		DB: db,
	}
	router := chi.NewRouter()
	router.Post("/signup", app.Signup)

	return router
}