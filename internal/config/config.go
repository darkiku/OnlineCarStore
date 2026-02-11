package config

import (
	"os"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	ServerPort   string
	JWTSecret    string
}

func Load() *Config {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "car_store"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	return &Config{
		MongoURI:     mongoURI,
		DatabaseName: dbName,
		ServerPort:   port,
		JWTSecret:    jwtSecret,
	}
}
