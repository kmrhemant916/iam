package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kmrhemant916/iam/controllers"
	"github.com/kmrhemant916/iam/helpers"
	"gopkg.in/yaml.v2"
)

type contextKey string

const (
	ConfigPath = "config/config.yaml"
	AuthHeader = "x-auth-token"
	ClaimsKey contextKey = "claims"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(AuthHeader)
		if tokenString == "" {
			response := map[string]interface{}{
				"message": "JWT is required",
			}
			helpers.SendResponse(w, response, http.StatusUnauthorized)
			return
		}
		jwtKey, err := GetJWTSecretKey()
		if err != nil {
			response := map[string]interface{}{
				"message": "Internal Server Error",
			}
			helpers.SendResponse(w, response, http.StatusInternalServerError)
			return
		}
		claims := &controllers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			response := map[string]interface{}{
				"message": "Unauthorized",
			}
			helpers.SendResponse(w, response, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetJWTSecretKey() ([]byte, error) {
	absPath, _ := helpers.GetAbsPath(ConfigPath)
	data, err := os.ReadFile(absPath)
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