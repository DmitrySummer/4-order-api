package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

type Config struct {
	UserEmail    string
	UserPassword string
	UserHost     string
	UserPort     string
	Db           DbConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файлов")
	}

	dbConfig := DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}

	return &Config{
		UserEmail:    os.Getenv("EMAIL"),
		UserPassword: os.Getenv("PASSWORD"),
		UserHost:     os.Getenv("HOST"),
		UserPort:     os.Getenv("PORT"),
		Db:           dbConfig,
	}
}
