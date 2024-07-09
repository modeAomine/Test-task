package Service

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
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

func GetUserByID(id int) (*Model.User, error) {
	var user Model.User
	err := DataBase.DB.QueryRow("SELECT id, username, hashed_password, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(u *Model.User) error {
	if u.Password != u.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hashedPassword)
	_, err = DataBase.DB.Exec("INSERT INTO users (username, hashed_password, role) VALUES ($1, $2, $3)", u.Username, u.HashedPassword, u.Role)
	return err
}

func UpdateUser(u *Model.User, currentPassword string) error {
	storedUser, err := GetUserByID(u.ID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(currentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	if u.Password != "" {
		if u.Password != u.ConfirmPassword {
			return errors.New("new password and confirmation do not match")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.HashedPassword = string(hashedPassword)
	} else {
		u.HashedPassword = storedUser.HashedPassword
	}

	_, err = DataBase.DB.Exec("UPDATE users SET username = $1, hashed_password = $2, role = $3 WHERE id = $4", u.Username, u.HashedPassword, u.Role, u.ID)
	return err
}

func DeleteUser(id int) error {
	_, err := DataBase.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
