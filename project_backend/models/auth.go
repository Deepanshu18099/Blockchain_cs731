package models

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password"`
	UserID   string `json:"userid"`
}

type Signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest12 struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type UserRequest23 struct {
	UserID string `json:"userid"`
	Email  string `json:"email"`
	Username   string `json:"username"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
}
