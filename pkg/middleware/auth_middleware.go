package middleware

import (
	"crimson-sunrise.site/pkg/persistence"
	"net/http"
	"strings"
	"context"
)

func AuthMiddleware(handler http.Handler) http.HandlerFunc {
	const userContextKey string = "user"
	return http.HandlerFunc(func (response http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(response,"Unauthorized", http.StatusUnauthorized)
			return
		}
		// now separate the Bearer part
		tokenString := strings.Split(authHeader,"Bearer ")[1]
		// validate the token with persistence layer
		user, err := persistence.VerifyToken(tokenString)
		if err != nil {
			http.Error(response,"Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(request.Context(),userContextKey, user)
		// continue the call chain
		handler.ServeHTTP(response,request.WithContext(ctx))
	})
}