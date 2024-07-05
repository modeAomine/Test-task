package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/http"
	"strconv"
	"tests/Config"
	"tests/DataBase"
	"tests/Router"
)

func main() {
	Config.LoadConfig()
	DataBase.Connect()

	m, err := migrate.New(
		"file://migrations",
		"postgres://"+Config.AppConfig.DB.Username+":"+Config.AppConfig.DBPassword+"@"+Config.AppConfig.DB.Host+":"+strconv.Itoa(Config.AppConfig.DB.Port)+"/"+Config.AppConfig.DB.DBName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	r := Router.AuthRouter()
	log.Fatal(http.ListenAndServe(":"+Config.AppConfig.Port, r))
}
