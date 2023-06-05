package main

import (
	"net/http"

	"github.com/kmrhemant916/iam/database"
	"github.com/kmrhemant916/iam/models"
	"github.com/kmrhemant916/iam/routes"
	"github.com/kmrhemant916/iam/utils"
)

const (
	Config = "config/config.yaml"
)

func main() {
	r := routes.SetupRoutes()
	var config utils.Config
	c, _:= config.ReadConf(Config)
	db,err := database.Connection()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Permission{})
	http.ListenAndServe(":"+c.Service.Port, r)
}