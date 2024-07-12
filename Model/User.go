package Model

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	HashedPassword string `json:"hashed_password"`
	Role           string `json:"role"`
}
