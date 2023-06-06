package main

import (
	"net/http"

	"github.com/kmrhemant916/iam/database"
	"github.com/kmrhemant916/iam/routes"
	"github.com/kmrhemant916/iam/utils"
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
	http.ListenAndServe(":"+c.Service.Port, r)
}