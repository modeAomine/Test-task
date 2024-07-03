package Config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	Port       string   `yaml:"port"`
	DB         DBConfig `yaml:"db"`
	DBPassword string
	JWTSecret  string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading config.yml file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &AppConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling config.yml: %v", err)
	}

	AppConfig.DBPassword = os.Getenv("DB_PASSWORD")
	if AppConfig.DBPassword == "" {
		log.Println("DB_PASSWORD is not set")
	}

	AppConfig.JWTSecret = os.Getenv("JWT_SECRET")
	if AppConfig.JWTSecret == "" {
		log.Println("JWT_SECRET is not set")
	}
}
