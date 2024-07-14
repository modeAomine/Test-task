package Validation

import (
	"regexp"
	"tests/DataBase"
)

type ValidationErrors map[string]string

func (ve ValidationErrors) Error() string {
	return "Validation errors occurred"
}

func ValidatePassword(password string) error {
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	specialCharRegex := regexp.MustCompile(`[!*]`)
	alphanumericWithSpecialRegex := regexp.MustCompile(`^[0-9a-zA-Z!*]*$`)
	lengthRegex := regexp.MustCompile(`^.{8,}$`)

	errors := ValidationErrors{}

	if !uppercaseRegex.MatchString(password) {
		errors["password"] = "Пароль должен содержать минимум 1 заглавную английскую букву!"
	}
	if !specialCharRegex.MatchString(password) {
		errors["password"] = "Пароль должен содержать минимум 1 специальный символ (!*)."
	}
	if !alphanumericWithSpecialRegex.MatchString(password) {
		errors["password"] = "Пароль должен содержать только английские буквы, цифры и специальные символы (!*)."
	}
	if !lengthRegex.MatchString(password) {
		errors["password"] = "Пароль должен состоять от 8 до 15 символов!"
	}

	if len(errors) > 0 {
		return errors
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
		return ValidationErrors{"username": "Пользователь с таким: " + username + " login уже существует!"}
	}

	return nil
}

func CheckUniqueEmail(email string) error {
	var count int
	err := DataBase.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ValidationErrors{"email": "Пользователь с таким: " + email + " email уже существует"}
	}

	return nil
}

func CheckUniquePhoneNumber(phone string) error {
	var count int
	err := DataBase.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phone = $1", phone).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ValidationErrors{"phone": "Пользователь с таким: " + phone + " уже занят!"}
	}

	return nil
}
