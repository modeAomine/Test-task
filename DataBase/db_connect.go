package DataBase

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	_ "log"
	"tests/Config"
)

var DB *sql.DB

func Connect() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",

		Config.AppConfig.DB.Host,
		Config.AppConfig.DB.Port,
		Config.AppConfig.DB.Username,
		Config.AppConfig.DBPassword,
		Config.AppConfig.DB.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
