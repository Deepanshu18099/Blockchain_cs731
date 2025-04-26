package models



type User struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRequest12 struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

type UserRequest23 struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	UserID string `json:"user_id"`
	Role  string `json:"role"`
}

type UserRequest32 struct {
	PubKey string `json:"pub_key"`
	PrivKey string `json:"priv_key"`
}


