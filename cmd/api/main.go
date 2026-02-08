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
	"github.com/teamserik/online-car-store/internal/middleware"
	"github.com/teamserik/online-car-store/internal/repository"
)

func main() {
	cfg := config.Load()

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

	carsCollection := database.GetCollection(client, cfg.DatabaseName, "cars")
	usersCollection := database.GetCollection(client, cfg.DatabaseName, "users")

	carRepo := repository.NewMongoCarRepository(carsCollection)
	userRepo := repository.NewMongoUserRepository(usersCollection)

	mux := http.NewServeMux()

	// Auth endpoints
	mux.HandleFunc("POST /api/auth/register", handler.Register(userRepo))
	mux.HandleFunc("POST /api/auth/login", handler.Login(userRepo))
	mux.HandleFunc("GET /api/auth/profile", middleware.AuthMiddleware(handler.GetProfile(userRepo)))

	// Car endpoints
	mux.HandleFunc("GET /api/cars", handler.ListCars(carRepo))
	mux.HandleFunc("GET /api/cars/", handler.GetCar(carRepo))
	mux.HandleFunc("POST /api/cars", middleware.AuthMiddleware(handler.CreateCar(carRepo)))
	mux.HandleFunc("PUT /api/cars/", middleware.AuthMiddleware(handler.UpdateCar(carRepo)))
	mux.HandleFunc("DELETE /api/cars/", middleware.AuthMiddleware(handler.DeleteCar(carRepo)))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	corsHandler := enableCORS(mux)

	addr := ":" + cfg.ServerPort
	fmt.Printf("Car Store API running on %s\n", addr)
	fmt.Printf("MongoDB Database: %s\n", cfg.DatabaseName)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
