package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/cors"
	"log"
	"net/http"
	"tests/Config"
	"tests/DataBase"
	"tests/Router"
)

func main() {
	Config.LoadConfig()
	DataBase.Connect()

	r := Router.MixRouter()
	r.Use(LoggingMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(r)

	log.Println("Server started on :" + Config.AppConfig.Port)
	log.Fatal(http.ListenAndServe(":"+Config.AppConfig.Port, handler))
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
