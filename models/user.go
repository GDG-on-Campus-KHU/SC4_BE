package models

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserData struct {
	Username string          `json:"username"`
	Supplies map[string]bool `json:"supplies"`
}
