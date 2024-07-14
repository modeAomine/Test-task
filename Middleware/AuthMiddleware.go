package Middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
	"tests/Utils"
)

var JwtKey = []byte("93821fsajkKFDS92")

type Claims struct {
	UserID int     `json:"user_id"`
	Role   string  `json:"role"`
	Phone  *string `json:"phone"`
	Email  *string `json:"email"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthMiddleware called for request: %s %s", r.Method, r.URL.Path)

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Token is valid, validating and extending if necessary")

		err = Utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/admin") && claims.Role != "admin" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
