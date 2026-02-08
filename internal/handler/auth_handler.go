package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/teamserik/online-car-store/internal/auth"
	"github.com/teamserik/online-car-store/internal/model"
	"github.com/teamserik/online-car-store/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(userRepo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.RegisterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.Email == "" || input.Password == "" || input.Username == "" {
			http.Error(w, "Email, username and password are required", http.StatusBadRequest)
			return
		}

		if len(input.Password) < 6 {
			http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		existingUser, _ := userRepo.FindByEmail(ctx, input.Email)
		if existingUser != nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}

		existingUser, _ = userRepo.FindByUsername(ctx, input.Username)
		if existingUser != nil {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}

		hashedPassword, err := auth.HashPassword(input.Password)
		if err != nil {
			http.Error(w, "Error processing password", http.StatusInternalServerError)
			return
		}

		user := &model.User{
			Username:  input.Username,
			Email:     input.Email,
			Password:  hashedPassword,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Phone:     input.Phone,
			Role:      "user",
		}

		if err := userRepo.Create(ctx, user); err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, user.Username, user.Role)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		response := model.AuthResponse{
			Token: token,
			User:  *user,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func Login(userRepo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.LoginInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Изменено: теперь проверяем username вместо email
		if input.Username == "" || input.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Изменено: ищем пользователя по username
		user, err := userRepo.FindByUsername(ctx, input.Username)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		if !auth.CheckPassword(input.Password, user.Password) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, user.Username, user.Role)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		response := model.AuthResponse{
			Token: token,
			User:  *user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetProfile(userRepo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := extractClaims(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func extractClaims(r *http.Request) (*auth.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("no authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("invalid authorization header")
	}

	return auth.ValidateToken(parts[1])
}
