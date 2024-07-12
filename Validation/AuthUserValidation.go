package Validation

import (
	"errors"
	"regexp"
)

func ValidateAuthUser(username, password string) error {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{4,10}$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("Длина имени пользователя должна составлять от 4 до 10 символов. Имя пользователя должно состоять только из английских букв и цифр!")
	}

	uppercaseRegex := regexp.MustCompile(`[A-Z].*[A-Z].*[A-Z]`)
	specialCharRegex := regexp.MustCompile(`[!*]`)
	alphanumericWithSpecialRegex := regexp.MustCompile(`^[0-9a-zA-Z!*]*$`)
	lengthRegex := regexp.MustCompile(`^.{8,15}$`)

	if !uppercaseRegex.MatchString(password) {
		return errors.New("Пароль должен содержать минимум 3 заглавные английские буквы!")
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
