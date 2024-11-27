package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	secretKey := []byte("2e5c28b16d282d501dcd65ac9132e9f8d053de116f14f411eb5bba846b90278b")

	claims := jwt.MapClaims{
		"user_id":  1,
		"username": "test",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error creating token:", err)
		return
	}

	fmt.Println("\n=== Test JWT Token ===")
	fmt.Println(tokenString)
	fmt.Println("\n=== Use this token in Authorization header ===")
	fmt.Println("Bearer", tokenString)
}
