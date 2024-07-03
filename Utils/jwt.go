package Utils

import (
	"github.com/golang-jwt/jwt/v4"
	"tests/Config"
	"tests/DataBase"
	"time"
)

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString([]byte(Config.AppConfig.JWTSecret))
}

func SaveTokenToDB(userID int, tokenString string, expiresAt time.Time) error {
	_, err := DataBase.DB.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, tokenString, expiresAt)
	return err
}
