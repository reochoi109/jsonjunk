package config

import (
	"fmt"
	"os"
)

type Config struct {
	MongoURI    string
	Port        string
	ServiceName string
	DBHost      string
	DBPort      string
	DBUsername  string
	DBPassword  string
	DBName      string
	ServiceMode string
	LogLevel    string
	Debug       bool
}

func LoadEnvConfig() *Config {
	return &Config{
		MongoURI:    buildMongoURI(),
		Port:        getEnv("PORT", "8080"),
		Debug:       getEnv("DEBUG", "true") == "true",
		ServiceName: getEnv("SERVICE_NAME", "pastebin-api"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "27017"),
		DBUsername:  getEnv("DB_USERNAME", ""),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "jsonjunk"),
		ServiceMode: getEnv("SERVICE_MODE", "dev"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	fmt.Println(fallback)
	return fallback
}

func buildMongoURI() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "27017")
	user := getEnv("DB_USERNAME", "")
	pass := getEnv("DB_PASSWORD", "")

	if user != "" && pass != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	}
	return fmt.Sprintf("mongodb://%s:%s", host, port)
}
