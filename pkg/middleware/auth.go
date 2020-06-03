package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/services"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("X-PackBack-Auth")
		if authHeader == "" {
			log.Println("Required header not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := &entities.User{}
		tokenService := services.NewTokenService(
			viper.GetString("token_service.secret_key"),
			viper.GetString("token_service.issuer"),
			viper.GetInt("token_service.minutes_valid"),
		)
		err := tokenService.Decode(authHeader, user)
		if err != nil {
			fmt.Printf("Error decoding token: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Println("User foud with email", " - ", user)

		ctx := context.WithValue(r.Context(), "User", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
