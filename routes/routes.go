package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kmrhemant916/iam/controllers"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)


func SetupRoutes(db *gorm.DB, conn *amqp.Connection) (*chi.Mux){
	app := &controllers.App{
		DB: db,
		Conn: conn,
	}
	router := chi.NewRouter()
	router.Post("/signup", app.Signup)

	return router
}