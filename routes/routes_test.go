package routes

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_routes(t *testing.T) {
	var registered = []struct{
		route string
		method string
	}{
		{"/signup", "POST"},
		
	}

	mux := SetupRoutes()



	for _, route := range registered {
		// check to see if the route exists
		if !routeExists(route.route, route.method, mux) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})
	return found
}