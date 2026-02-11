package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	favoritesCollection := database.GetCollection(client, cfg.DatabaseName, "favorites")
	reviewsCollection := database.GetCollection(client, cfg.DatabaseName, "reviews")

	carRepo := repository.NewMongoCarRepository(carsCollection)
	userRepo := repository.NewMongoUserRepository(usersCollection)
	favoriteRepo := repository.NewMongoFavoriteRepository(favoritesCollection, carsCollection)
	reviewRepo := repository.NewMongoReviewRepository(reviewsCollection)

	mux := http.NewServeMux()

	// Auth endpoints
	mux.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.Register(userRepo)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.Login(userRepo)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/auth/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			middleware.AuthMiddleware(handler.GetProfile(userRepo))(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Car endpoints
	mux.HandleFunc("/api/cars", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cars" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handler.ListCars(carRepo)(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(handler.CreateCar(carRepo))(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cars/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/cars" || r.URL.Path == "/api/cars/" {
			http.Redirect(w, r, "/api/cars", http.StatusMovedPermanently)
			return
		}

		// Извлекаем ID из пути
		path := strings.TrimPrefix(r.URL.Path, "/api/cars/")

		// Проверяем если это запрос отзывов
		if strings.HasSuffix(path, "/reviews") {
			carID := strings.TrimSuffix(path, "/reviews")

			switch r.Method {
			case http.MethodPost:
				// Добавляем car_id в контекст запроса для handler
				r.URL.RawQuery = "car_id=" + carID
				middleware.AuthMiddleware(handler.CreateReview(reviewRepo, userRepo))(w, r)
			case http.MethodGet:
				r.URL.RawQuery = "car_id=" + carID
				handler.GetCarReviews(reviewRepo)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Обычные операции с машинами
		switch r.Method {
		case http.MethodGet:
			handler.GetCar(carRepo)(w, r)
		case http.MethodPut:
			middleware.AuthMiddleware(handler.UpdateCar(carRepo))(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(handler.DeleteCar(carRepo))(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Favorites endpoints
	mux.HandleFunc("/api/favorites", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/favorites" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			middleware.AuthMiddleware(handler.GetFavorites(favoriteRepo))(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(handler.AddToFavorites(favoriteRepo))(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/favorites/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			middleware.AuthMiddleware(handler.GetFavoritesCount(favoriteRepo))(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/favorites/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/favorites" || r.URL.Path == "/api/favorites/" {
			http.Redirect(w, r, "/api/favorites", http.StatusMovedPermanently)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/favorites/count") {
			return // Уже обработано выше
		}

		if r.Method == http.MethodDelete {
			middleware.AuthMiddleware(handler.RemoveFromFavorites(favoriteRepo))(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Reviews endpoints
	mux.HandleFunc("/api/reviews", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/reviews" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handler.GetCarReviews(reviewRepo)(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(handler.CreateReview(reviewRepo, userRepo))(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/reviews/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/reviews" || r.URL.Path == "/api/reviews/" {
			http.Redirect(w, r, "/api/reviews", http.StatusMovedPermanently)
			return
		}

		switch r.Method {
		case http.MethodPut:
			middleware.AuthMiddleware(handler.UpdateReview(reviewRepo))(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(handler.DeleteReview(reviewRepo))(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Static files
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
