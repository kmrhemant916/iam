package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kmrhemant916/iam/controllers"
	"github.com/kmrhemant916/iam/middlewares"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)


func SetupRoutes(db *gorm.DB, conn *amqp.Connection) (*chi.Mux){
	var jwtKey []byte
	jwtKey, _ = middlewares.GetJWTSecretKey()
	app := controllers.NewApp(db, conn, jwtKey)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/signup", app.Signup)
	router.Post("/signin", app.Signin)
	router.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Get("/roles", app.GetRoles)
		r.Get("/api/users/{id}/profile", app.GetUserProfile)
		r.Post("/api/users/invite", app.InviteUser)
	})

	return router
}