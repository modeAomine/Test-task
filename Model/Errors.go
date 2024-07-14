package Model

import "errors"

var (
	ErrWardrobeNotFound    = errors.New("Шкаф не найден")
	ErrInvalidWardrobeData = errors.New("Не верные данные шкафа")
	ErrUserNotFound        = errors.New("пользователь не найден")
	ErrInvalidUserData     = errors.New("неверные данные пользователя")
	ErrInvalidPassword     = errors.New("Не верный пароль")
)
