package controllers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
	Conn *amqp.Connection
	JWTKey []byte
}

func NewApp(db *gorm.DB, conn *amqp.Connection, jwtKey []byte) (*App) {
	app := &App{
		DB:          db,
		Conn:        conn,
		JWTKey:      jwtKey,
	}
	return app
}