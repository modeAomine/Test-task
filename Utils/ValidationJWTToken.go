package Utils

import (
	"database/sql"
	"errors"
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

	if time.Now().After(expiresAt) {
		return errors.New("token expired")
	}

	return nil
}
