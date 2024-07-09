package Service

import (
	"database/sql"
	"tests/DataBase"
	"tests/Model"
	"time"
)

func GetUserByUsername(username string) (*Model.User, error) {
	var user Model.User
	err := DataBase.DB.QueryRow("SELECT id, username, hashed_password, role FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetActiveTokenByUserID(userID int) (string, error) {
	var token string
	var expiresAt time.Time

	err := DataBase.DB.QueryRow("SELECT token, expires_at FROM tokens WHERE user_id = $1 AND expires_at > $2", userID, time.Now()).Scan(&token, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return token, nil
}
