package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/teamserik/online-car-store/internal/config"
	"github.com/teamserik/online-car-store/internal/database"
	"github.com/teamserik/online-car-store/internal/handler"
	"github.com/teamserik/online-car-store/internal/repository"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к MongoDB
	client, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Получение коллекции
	carsCollection := database.GetCollection(client, cfg.DatabaseName, "cars")
	carRepo := repository.NewMongoCarRepository(carsCollection)

	// Настройка роутов
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("POST /api/cars", handler.CreateCar(carRepo))
	mux.HandleFunc("GET /api/cars", handler.ListCars(carRepo))
	mux.HandleFunc("GET /api/cars/", handler.GetCar(carRepo))
	mux.HandleFunc("PUT /api/cars/", handler.UpdateCar(carRepo))
	mux.HandleFunc("DELETE /api/cars/", handler.DeleteCar(carRepo))

	// Статические файлы для фронтенда
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	// CORS middleware
	corsHandler := enableCORS(mux)

	// Запуск сервера
	addr := ":" + cfg.ServerPort
	fmt.Printf("Car Store API running on %s\n", addr)
	fmt.Printf("MongoDB Database: %s\n", cfg.DatabaseName)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
