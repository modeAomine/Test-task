package Model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"tests/DataBase"
)

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	HashedPassword string `json:"hashed_password"`
}

func (u *User) Create() error {
	if u.Password != u.HashedPassword {
		return errors.New("passwords do not match")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = DataBase.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, string(hashedPassword))
	return err
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DataBase.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.HashedPassword)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
