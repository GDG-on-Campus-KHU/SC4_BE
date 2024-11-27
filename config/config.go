package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	DB        *DBConfig
	JWTSecret []byte
}

func GetConfig() *Config {
	return &Config{
		DB:        GetDBConfig(),
		JWTSecret: []byte(getEnv("JWT_SECRET", "2e5c28b16d282d501dcd65ac9132e9f8d053de116f14f411eb5bba846b90278b")),
	}
}

func GetDBConfig() *DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "gdgdb"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
