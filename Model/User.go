package Model

type User struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	HashedPassword  string `json:"hashed_password"`
	Role            string `json:"role"`
}
