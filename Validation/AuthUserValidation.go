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

	passwordRegex := regexp.MustCompile(`^(?=(.*[A-Z]){3})(?=.*[!*])(?=.*[0-9a-zA-Z]).{8,12}$`)
	if !passwordRegex.MatchString(password) {
		return errors.New("Пароль пользователя должен состоять от 8 до 12 символов, также должен содержать минимум 3 заглавные английские буквы и 1 специальный символ (!*). Пароль пользователя должен содержать только английские буквы и цифры!")
	}

	return nil
}
