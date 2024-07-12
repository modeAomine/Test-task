package Service

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"tests/DataBase"
	"tests/Model"
	"time"
)

func RegisterUser(user *Model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedHashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.HashedPassword = string(hashedHashedPassword)

	err = CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(user *Model.User) error {
	_, err := DataBase.DB.Exec("INSERT INTO users (username, password, hashed_password, role) VALUES ($1, $2, $3, $4)", user.Username, user.Password, user.HashedPassword, user.Role)
	return err
}

func GetUserByUsername(username string) (*Model.User, error) {
	var user Model.User
	err := DataBase.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role)
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

func GetAllUsers() ([]Model.User, error) {
	var users []Model.User
	rows, err := DataBase.DB.Query("SELECT id, username, password, hashed_password, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user Model.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.HashedPassword,
			&user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
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

func UpdateUserByAdmin(user *Model.User) error {
	storedUser, err := GetUserByID(user.ID)
	if err != nil {
		return err
	}
	if storedUser == nil {
		return errors.New("User not found")
	}

	storedUser.Username = user.Username
	storedUser.Role = user.Role

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		storedUser.Password = string(hashedPassword)
	}

	_, err = DataBase.DB.Exec(`
		UPDATE users 
		SET username = $1, hashed_password = $2, role = $3
		WHERE id = $4`,
		storedUser.Username,
		storedUser.HashedPassword,
		storedUser.Role,
		storedUser.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserProfile(user *Model.User, currentPassword string) error {
	storedUser, err := GetUserByID(user.ID)
	if err != nil {
		return err
	}
	if storedUser == nil {
		return errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(currentPassword))
	if err != nil {
		return errors.New("Invalid password")
	}

	storedUser.Username = user.Username

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		storedUser.HashedPassword = string(hashedPassword)
	}

	NewPasswordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DataBase.DB.Exec(`
		UPDATE users 
		SET username = $1, password = $2, hashed_password = $3
		WHERE id = $4`,
		storedUser.Username,
		string(NewPasswordHash),
		storedUser.HashedPassword,
		storedUser.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id int) error {
	_, err := DataBase.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
