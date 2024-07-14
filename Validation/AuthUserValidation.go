package Validation

import (
	"errors"
	"regexp"
	"tests/DataBase"
)

func ValidateAuthUser(username, password string) error {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("Длина login'a пользователя должна составлять от 4 до 15 символов. Имя пользователя должно состоять только из английских букв и цифр!")
	}

	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	specialCharRegex := regexp.MustCompile(`[!*]`)
	alphanumericWithSpecialRegex := regexp.MustCompile(`^[0-9a-zA-Z!*]*$`)
	lengthRegex := regexp.MustCompile(`^.{8,}$`)

	if !uppercaseRegex.MatchString(password) {
		return errors.New("Пароль должен содержать минимум 1 заглавную английскую букву!")
	}
	if !specialCharRegex.MatchString(password) {
		return errors.New("Пароль должен содержать минимум 1 специальный символ (!*).")
	}
	if !alphanumericWithSpecialRegex.MatchString(password) {
		return errors.New("Пароль должен содержать только английские буквы, цифры и специальные символы (!*).")
	}
	if !lengthRegex.MatchString(password) {
		return errors.New("Пароль должен состоять от 8 до 15 символов!")
	}

	return nil
}

func CheckUniqueUsername(username string) error {
	var count int
	err := DataBase.DB.QueryRow("select COUNT(*) from users where username = $1", username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Пользователь с таким: " + username + " login уже существует!")
	}

	return nil
}

func CheckUniqueEmailAndPhone(email string, phone string) error {
	var count int
	err := DataBase.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Пользователь с таким: " + email + " email уже существует")
	}

	err = DataBase.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phone = $1", phone).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("Пользователь с таким: " + phone + " номером уже существует!")
	}
	return nil
}
