package Controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tests/Model"
	"tests/Service"
)

func UpdateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var updateRequest struct {
		User Model.User `json:"user"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	user := updateRequest.User
	user.ID = id

	err = Service.UpdateUserByAdmin(&user)
	if err != nil {
		switch err {
		case Model.ErrUserNotFound:
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		case Model.ErrInvalidUserData:
			http.Error(w, "Неверные данные пользователя", http.StatusBadRequest)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"message": "Пользователь успешно обновлен", "user": user}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = Service.DeleteUser(id)
	if err != nil {
		switch err {
		case Model.ErrUserNotFound:
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"message": "Пользователь успешно удален"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user Model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	err = Service.CreateUser(&user)
	if err != nil {
		switch err {
		case Model.ErrInvalidUserData:
			http.Error(w, "Неверные данные пользователя", http.StatusBadRequest)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := Service.GetAllUsers()
	if err != nil {
		http.Error(w, "Не удалось получить пользователей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"users": users}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Не удалось закодировать ответ", http.StatusInternalServerError)
	}
}

func GetAllWardrobe(w http.ResponseWriter, r *http.Request) {
	/*	if !Service.CheckTokenExpiration(w, r) {
		return
	}*/
	wardrobe, err := Service.GetAllWardrobe()
	if err != nil {
		http.Error(w, "Не удалось получить шкафы", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"wardrobe": wardrobe}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Не удалось закодировать ответ", http.StatusInternalServerError)
		return
	}
}
