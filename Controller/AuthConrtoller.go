package Controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"tests/Model"
	"tests/Response"
	"tests/Service"
	"tests/Utils"
	"tests/Validation"
	"time"
)

type RegistrationRequest struct {
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	HashedPassword string `json:"hashed_password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := Validation.ValidationErrors{}

	if err := Validation.CheckUniqueUsername(req.Username); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := Validation.CheckUniqueEmail(req.Email); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := Validation.CheckUniquePhoneNumber(req.Phone); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := Validation.CheckUniqueUsername(req.Username); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := Validation.ValidatePassword(req.Password); err != nil {
		if verr, ok := err.(Validation.ValidationErrors); ok {
			for k, v := range verr {
				errors[k] = v
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if req.Password != req.HashedPassword {
		errors["confirmPassword"] = "Пароль и подтверждение пароля не совпадают!"
	}

	if req.FullName == "" {
		errors["fullName"] = "Имя пользователя не может быть пустым!"
	}

	if req.Email == "" {
		errors["email"] = "Почта пользователя не может быть пустой!"
	}

	if req.Phone == "" {
		errors["phone"] = "Номер телефона не может быть пустым!"
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errors})
		return
	}

	user := Model.User{
		Username:       req.Username,
		FullName:       sql.NullString{String: req.FullName, Valid: req.FullName != ""},
		Email:          sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Password:       req.Password,
		HashedPassword: req.HashedPassword,
		Role:           "user",
	}

	err = Service.RegisterUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": user.Username + " успешно зарегистрирован!",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Response.SendError(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	storedUser, err := Service.GetUserByUsername(req.Username)
	if err != nil {
		Response.SendError(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(req.Password))
	if err != nil {
		Response.SendError(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	activeToken, err := Service.GetActiveTokenByUserID(storedUser.ID)
	if err != nil {
		Response.SendError(w, "Ошибка при проверке токена", http.StatusInternalServerError)
		return
	}

	if activeToken != "" {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message": "Вы уже авторизованы!",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := Utils.GenerateJWT(storedUser.ID, storedUser.Username, storedUser.Role, storedUser.Email, storedUser.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Generated Token: %s\n", token)

	expiresAt := time.Now().Add(time.Hour * 12)
	err = Utils.SaveTokensToDB(storedUser.ID, token, expiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		Response.SendError(w, "Требуется заголовок авторизации", http.StatusUnauthorized)
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	log.Printf("Token to invalidate: %s", tokenString)

	err := Service.InvalidToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Успешный выход из системы!",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
