package Service

import (
	"errors"
	"tests/DataBase"
)

func InvalidToken(token string) error {
	var count int
	err := DataBase.DB.QueryRow("SELECT COUNT(*) FROM tokens WHERE token = $1", token).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Токен не найден")
	}

	_, err = DataBase.DB.Exec("DELETE FROM tokens WHERE token = $1", token)
	return err
}

/*func CheckTokenExpiration(w http.ResponseWriter, r *http.Request) bool {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Требуется заголовок авторизации", http.StatusUnauthorized)
		return false
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "JWT_SECRET не установлен", http.StatusInternalServerError)
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				Controller.Logout(w, r)
				return false
			}
		}
	} else {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return false
	}

	return true
}*/
