package Controller

import (
	"database/sql"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tests/DataBase"
	"tests/Middleware"
	"tests/Model"
	"tests/Service"
)

type UserProfile struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*Middleware.Claims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userID := claims.UserID

	row := DataBase.DB.QueryRow("SELECT username, email, full_name FROM users WHERE id = $1", userID)

	var userProfile UserProfile
	err := row.Scan(&userProfile.Username, &userProfile.Email, &userProfile.FullName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	if !Service.CheckTokenExpiration(w, r) {
		return
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Требуется заголовок авторизации", http.StatusUnauthorized)
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "JWT_SECRET не установлен", http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFromToken := int(claims["user_id"].(float64))

		if userIDFromToken != id {
			http.Error(w, "Вы можете обновить только свой профиль", http.StatusForbidden)
			return
		}
	} else {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	var updateRequest struct {
		CurrentPassword string     `json:"current_password"`
		User            Model.User `json:"user"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	user := updateRequest.User
	user.ID = id

	err = Service.UpdateUserProfile(&user, updateRequest.CurrentPassword)
	if err != nil {
		switch err {
		case Model.ErrUserNotFound:
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		case Model.ErrInvalidPassword:
			http.Error(w, "Неверный текущий пароль", http.StatusBadRequest)
		case Model.ErrInvalidUserData:
			http.Error(w, "Неверные данные пользователя", http.StatusBadRequest)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"message": "Профиль пользователя успешно изменен", "user": user}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
