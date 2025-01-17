package Controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tests/Model"
	"tests/Response"
	"tests/Service"
	"tests/Validation"
)

func UpdateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	if !Service.CheckTokenExpiration(w, r) {
		return
	}
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

	errors := make(map[string]string)

	err = Service.UpdateUserByAdmin(&user)
	if err != nil {
		switch err {
		case Model.ErrUserNotFound:
			errors["id"] = "Пользователь не найден"
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
	if !Service.CheckTokenExpiration(w, r) {
		return
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	errors := make(map[string]string)

	err = Service.DeleteUser(id)
	if err != nil {
		switch err {
		case Model.ErrUserNotFound:
			errors["id"] = "Пользователь не найден"
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

	errors := Validation.ValidationErrors{}

	if err := Validation.CheckUniqueEmail(user.Email.String); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := Validation.CheckUniquePhoneNumber(user.Phone.String); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := Validation.CheckUniqueUsername(user.Username); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := Validation.ValidatePassword(user.Password); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if user.Username == "" {
		errors["username"] = "Login пользователя не может быть пустым"
	}

	if user.Password != user.HashedPassword {
		errors["confirmPassword"] = "Пароль и подтверждение пароля не совпадают!"
	}

	if user.FullName.String == "" {
		errors["fullName"] = "Имя пользователя не может быть пустым!"
	}

	if user.Email.String == "" {
		errors["email"] = "Почта пользователя не может быть пустой!"
	}

	if user.Phone.String == "" {
		errors["phone"] = "Номер телефона не может быть пустым!"
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response.ErrorResponse{Errors: errors})
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
	if !Service.CheckTokenExpiration(w, r) {
		return
	}
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
