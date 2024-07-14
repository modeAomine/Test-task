package Utils

import (
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"tests/Config"
	"tests/DataBase"
	"time"
)

func GenerateJWT(userID int, username string, role string, email sql.NullString, phone sql.NullString) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}

	if email.Valid {
		claims["email"] = email.String
	}

	if phone.Valid {
		claims["phone"] = phone.String
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := token.SignedString([]byte(Config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func SaveTokensToDB(userID int, jwtToken string, expiresAt time.Time) error {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	expiresAt = expiresAt.In(location)

	_, err = DataBase.DB.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, jwtToken, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func ExtendTokenLifetime(token string, newExpiresAt time.Time) error {
	log.Printf("Extending token lifetime for token: %s to: %v", token, newExpiresAt)
	_, err := DataBase.DB.Exec("UPDATE tokens SET expires_at = $1 WHERE token = $2", newExpiresAt, token)
	if err != nil {
		log.Printf("Failed to extend token lifetime: %v", err)
		return err
	}
	log.Printf("Token lifetime successfully extended")
	return nil
}
