package main

import (
	"net/http"

	"github.com/kmrhemant916/iam/database"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/rabbitmq"
	"github.com/kmrhemant916/iam/routes"
)

const (
	Config = "config/config.yaml"
)

func main() {
	var config helpers.Config
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
	rabbitmqConfig := c.Rabbitmq
	conn, err := rabbitmq.Connection(rabbitmqConfig.Username, rabbitmqConfig.Password, rabbitmqConfig.Host, rabbitmqConfig.Port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	r := routes.SetupRoutes(db, conn)
	helpers.InitialiseServices(db)
	helpers.InitialiseAuthorization(db, c.Roles, c.Permissions, c.RolePermissions)
	http.ListenAndServe(":"+c.Service.Port, r)
}