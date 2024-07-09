package Utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"tests/Config"
	"tests/DataBase"
	"time"
)

func GenerateJWT(userID int, username string, role string) (string, error) {
	fmt.Println("Role: ", role+"Username: ", username+"UserID: ", userID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString([]byte(Config.AppConfig.JWTSecret))
}

func SaveTokenToDB(userID int, tokenString string, expiresAt time.Time) error {
	_, err := DataBase.DB.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, tokenString, expiresAt)
	return err
}

func DecodeJWT(tokenString string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.AppConfig.JWTSecret), nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("User ID:", claims["user_id"])
		fmt.Println("Role:", claims["role"])
	} else {
		fmt.Println("Invalid token")
	}
}
