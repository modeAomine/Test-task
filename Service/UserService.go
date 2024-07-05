package Service

import (
	"tests/DataBase"
	"tests/Model"
)

func GetUserByUsername(username string) (*Model.User, error) {
	var user Model.User
	err := DataBase.DB.QueryRow("SELECT id, username, hashed_password, role FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
