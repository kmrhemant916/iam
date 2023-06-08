package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/kmrhemant916/iam/database"
	"github.com/kmrhemant916/iam/rabbitmq"
	"github.com/kmrhemant916/iam/routes"
	"github.com/kmrhemant916/iam/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Config = "config/config.yaml"
)

func main() {
	var config utils.Config
	c, err:= config.ReadConf(Config)
    if err != nil {
        panic(err)
    }
	dbConfig := c.Database
	db,err := database.Connection(dbConfig.Host, dbConfig.Name, dbConfig.Password, dbConfig.Port, dbConfig.Username)
	if err != nil {
		panic(err)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()
	r := routes.SetupRoutes(db)
	utils.InitialiseServices(db)
	rabbitmqConfig := c.Rabbitmq
	conn, err := rabbitmq.Connection(rabbitmqConfig.Username, rabbitmqConfig.Password, rabbitmqConfig.Host, rabbitmqConfig.Port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"mail", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
	  "",     // exchange
	  q.Name, // routing key
	  false,  // mandatory
	  false,  // immediate
	  amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	  })
	if err != nil {
		panic(err)
	}
	log.Printf(" [x] Sent %s\n", body)
	http.ListenAndServe(":"+c.Service.Port, r)
}