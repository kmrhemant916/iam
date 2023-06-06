package routes

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kmrhemant916/iam/database"
	"github.com/kmrhemant916/iam/utils"
)

const (
	Config = "../config/config.yaml"
)

func TestAppRoutes(t *testing.T) {
	registered := []struct {
		route  string
		method string
	}{
		{"/signup", "POST"},
	}
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
	mux := SetupRoutes(db)
	for _, route := range registered {
		// Check to see if the route exists
		if !routeExists(route.route, route.method, mux) {
			t.Errorf("route %s %s is not registered", route.method, route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		if method == testMethod && route == testRoute {
			found = true
		}
		return nil
	})
	return found
}