package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"os"
	"strings"
)

type MiddlewareFunc func(http.Handler) http.Handler

func ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//get the authorization header which contains a bearer token
		authorizationHeader := r.Header.Get("Authorization")
		log.Print(authorizationHeader)
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					//jwt-sample is the JWT secret this should be changed to a secure secret
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if err != nil {
					log.Print(err)
					_ = json.NewEncoder(w).Encode(err.Error())
					return
				}
				if token.Valid {
					context.Set(r, "decoded", token.Claims)
					next.ServeHTTP(w, r)
				} else {
					_ = json.NewEncoder(w).Encode("Invalid authorization token")
				}
			} else {
				_ = json.NewEncoder(w).Encode("Invalid authorization token")
			}
		} else {
			_ = json.NewEncoder(w).Encode("An authorization header is required")
		}
	})
}
