package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	RestID    int64     `json:"rest_id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	PhoneNum  string    `json:"phone_num"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	Username string          `json:"username"`
	Supplies map[string]bool `json:"supplies"`
}
