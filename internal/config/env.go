package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("could not load .env file", err.Error())
	}

	return Config{
		Port: getEnv("PORT", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return ":" + value
	}

	return fallback
}
