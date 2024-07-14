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

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = Validation.CheckUniqueUsername(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = Validation.CheckUniqueEmailAndPhone(req.Email, req.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = Validation.ValidateAuthUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password != req.HashedPassword {
		http.Error(w, "Пароль и подтверждение пароля не совпадают!", http.StatusBadRequest)
		return
	}

	if req.FullName == "" {
		http.Error(w, "Имя пользователя не может быть пустым!", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Почта пользователя не может быть пустой!", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		http.Error(w, "Номер телефона не может быть пустым!", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedUser, err := Service.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	activeToken, err := Service.GetActiveTokenByUserID(storedUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Требуется заголовок авторизации", http.StatusUnauthorized)
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
