package Service

import (
	"errors"
	"tests/DataBase"
)

func InvalidToken(token string) error {
	var count int
	err := DataBase.DB.QueryRow("SELECT COUNT(*) FROM tokens WHERE token = $1", token).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Токен не найден")
	}

	_, err = DataBase.DB.Exec("DELETE FROM tokens WHERE token = $1", token)
	return err
}
