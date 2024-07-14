package Utils

import (
	"database/sql"
	"errors"
	"log"
	"tests/DataBase"
	"time"
)

func ValidateToken(tokenString string) error {
	var expiresAt time.Time
	err := DataBase.DB.QueryRow("SELECT expires_at FROM tokens WHERE token = $1", tokenString).Scan(&expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("token not found")
		}
		return err
	}

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Printf("Failed to load location: %v", err)
		return err
	}

	expiresAt = expiresAt.In(location)

	log.Printf("Current time: %v", time.Now().In(location))
	log.Printf("Token expires at: %v", expiresAt)
	log.Printf("Time difference: %v", expiresAt.Sub(time.Now().In(location)))

	if expiresAt.Sub(time.Now().In(location)) < 3*time.Hour+30*time.Minute {
		newExpiresAt := time.Now().Add(time.Hour * 12)
		log.Printf("Extending token lifetime to: %v", newExpiresAt)
		err := ExtendTokenLifetime(tokenString, newExpiresAt)
		if err != nil {
			log.Printf("Failed to extend token lifetime: %v", err)
			return err
		}
	} else {
		log.Printf("Token lifetime does not need extension")
	}

	return nil
}
