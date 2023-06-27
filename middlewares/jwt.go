package middlewares

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kmrhemant916/iam/controllers"
	"github.com/kmrhemant916/iam/helpers"
	"gopkg.in/yaml.v2"
)
const (
	ConfigPath = "config/config.yaml"
	AuthHeader = "x-auth-token"
)
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(AuthHeader)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		jwtKey, err := getJWTSecretKey()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		claims := &controllers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getJWTSecretKey() ([]byte, error) {
	absPath, _ := helpers.GetAbsPath(ConfigPath)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	jwtSecret, ok := config["jwt_key"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid JWT secret key")
	}
	return []byte(jwtSecret), nil
}