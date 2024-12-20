package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

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

	// c:= cors.New(cors.Options{
	// 	AllowedOrigins: []string{
	// 		"https://packmate-three.vercel.app",
	// 		"http://localhost:5173",
	// 		"http://localhost:5500",
	// 	},
	// 	AllowCredentials: true,
	r := mux.NewRouter()
	c := cors.AllowAll()

	userService := &services.UserService{}
	userHandler := handlers.NewUserHandler(userService)

	suppliesService := services.NewSuppliesService()
	suppliesHandler := handlers.NewSuppliesHandler(suppliesService, cfg)

	//api

	r.HandleFunc("/api/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/api/register", userHandler.CreateUser).Methods("POST")

	//JWT 적용후
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(handlers.AuthMiddleware)
	//protected.Use(handlers.CORSMiddleware(allowedOrigins))

	protected.HandleFunc("/user", suppliesHandler.GetSupplies).Methods("GET")
	protected.HandleFunc("/supplies", suppliesHandler.SaveSupplies).Methods("POST")
	protected.HandleFunc("/supplies", suppliesHandler.UpdateSupplies).Methods("PUT")
	//JWT 적용전
	//r.HandleFunc("/api/v1/user", suppliesHandler.GetSupplies).Methods("GET")
	// r.HandleFunc("/api/v1/supplies", suppliesHandler.SaveSupplies).Methods("POST")
	// r.HandleFunc("/api/v1/supplies", suppliesHandler.UpdateSupplies).Methods("PUT")

	log.Println("Server starting at :8888")
	log.Fatal(http.ListenAndServe(":8888", c.Handler(r)))
}
