package Controller

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"tests/Model"
	"tests/Service"
	"tests/Utils"
	"time"
)

type RegistrationRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	HashedPassword string `json:"hashed_password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Registration(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := Model.User{
		Username:       req.Username,
		Password:       req.Password,
		HashedPassword: req.HashedPassword,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	err = Service.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": user.Username + " " + "успешно зарегестрировался!",
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

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(req.Password))
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
			"message": "Вы уже авторезированы!",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := Utils.GenerateJWT(storedUser.ID, storedUser.Username, storedUser.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(time.Hour * 1)
	err = Utils.SaveTokenToDB(storedUser.ID, token, expiresAt)
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
