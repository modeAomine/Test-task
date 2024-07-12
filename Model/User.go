package Model

import "database/sql"

type User struct {
	ID             int            `json:"id"`
	Username       string         `json:"username"`
	FullName       sql.NullString `json:"full_name"`
	Email          sql.NullString `json:"email"`
	Phone          sql.NullString `json:"phone"`
	Password       string         `json:"password"`
	HashedPassword string         `json:"hashed_password"`
	Role           string         `json:"role"`
}
