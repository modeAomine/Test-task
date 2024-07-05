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
	Role           string `json:"role"`
}

func (u *User) Create() error {
	if u.Password != u.HashedPassword {
		return errors.New("passwords do not match")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hashedPassword)
	_, err = DataBase.DB.Exec("INSERT INTO users (username, password, hashed_password) VALUES ($1, $2, $3)", u.Username, string(hashedPassword), u.HashedPassword)
	return err
}
