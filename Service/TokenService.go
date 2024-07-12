package Service

import "tests/DataBase"

func InvalidToken(token string) error {
	_, err := DataBase.DB.Exec("DELETE FROM Token WHERE token = $1", token)
	return err
}
