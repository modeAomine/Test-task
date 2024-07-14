package Model

import "errors"

var (
	ErrWardrobeNotFound    = errors.New("Шкаф не найден")
	ErrInvalidWardrobeData = errors.New("Не верные данные шкафа")
)
