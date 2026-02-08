package config

import (
	"os"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	ServerPort   string
}

func Load() *Config {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb+srv://banichban186_db_user:Zhmblkhn2627@onlinecarstore.6ddaifq.mongodb.net/?"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "car_store"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	return &Config{
		MongoURI:     mongoURI,
		DatabaseName: dbName,
		ServerPort:   port,
	}
}
