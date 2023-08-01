package middlewares

import (
	"net/http"

	"github.com/kmrhemant916/iam/controllers"
	"github.com/kmrhemant916/iam/global"
	"github.com/kmrhemant916/iam/helpers"
)

const (
	GroupSendInvite = "Administrator"
)

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(controllers.ClaimsKey).(*controllers.Claims)
		if !ok {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		authorized := false
		for _, group := range claims.Groups {
			if group == GroupSendInvite {
				authorized = true
				break
			}
		}
		if !authorized {
			response := map[string]interface{}{
				"message": "Unauthorized",
			}
			helpers.SendResponse(w, response, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

