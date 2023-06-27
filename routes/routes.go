package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kmrhemant916/iam/controllers"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/middlewares"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)


func SetupRoutes(db *gorm.DB, conn *amqp.Connection, c *helpers.Config) (*chi.Mux){
	app := &controllers.App{
		DB: db,
		Conn: conn,
		JWTKey: c.JWTKey,
	}
	router := chi.NewRouter()
	router.Post("/signup", app.Signup)
	router.Post("/signin", app.Signin)
	router.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Get("/api", app.Roles)
	})

	return router
}