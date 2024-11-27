package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GDG-on-Campus-KHU/SC4_BE/config"
	"github.com/GDG-on-Campus-KHU/SC4_BE/db"
	"github.com/GDG-on-Campus-KHU/SC4_BE/handlers"
	"github.com/GDG-on-Campus-KHU/SC4_BE/services"
)

func main() {
	cfg := config.GetConfig()

	err := db.InitDB(cfg.DB)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.DB.Close()

	r := mux.NewRouter()

	userService := &services.UserService{}
	userHandler := handlers.NewUserHandler(userService)

	suppliesService := services.NewSuppliesService()
	suppliesHandler := handlers.NewSuppliesHandler(suppliesService, cfg)

	//api
	r.HandleFunc("/api/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/api/register", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/user", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/user/{id}", userHandler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/api/v1/user", suppliesHandler.GetSupplies).Methods("GET")
	r.HandleFunc("/api/v1/supplies", suppliesHandler.SaveSupplies).Methods("POST")
	r.HandleFunc("/api/v1/supplies", suppliesHandler.UpdateSupplies).Methods("PUT")

	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
