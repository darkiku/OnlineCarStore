package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/teamserik/online-car-store/internal/model"
	"github.com/teamserik/online-car-store/internal/repository"
)

var viewCount int64 // атомарный счётчик просмотров (для горутины)

func CreateCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var input model.CreateCarInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// минимальная проверка
		if input.Make == "" || input.Model == "" || input.Year < 1900 || input.Price <= 0 {
			respondError(w, http.StatusBadRequest, "Required fields missing or invalid")
			return
		}

		car := &model.Car{
			Make:        input.Make,
			Model:       input.Model,
			Year:        input.Year,
			Price:       input.Price,
			Mileage:     input.Mileage,
			BodyType:    input.BodyType,
			Description: input.Description,
		}

		repo.Create(car)
		respondJSON(w, http.StatusCreated, car)
	}
}

func ListCars(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		cars := repo.List()
		respondJSON(w, http.StatusOK, cars)
	}
}

func GetCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/cars/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid car ID")
			return
		}

		car, found := repo.GetByID(id)
		if !found {
			respondError(w, http.StatusNotFound, "Car not found")
			return
		}

		atomic.AddInt64(&viewCount, 1) // для статистики в фоне
		respondJSON(w, http.StatusOK, car)
	}
}

func UpdateCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/cars/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid car ID")
			return
		}

		var input model.UpdateCarInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if !repo.Update(id, input) {
			respondError(w, http.StatusNotFound, "Car not found")
			return
		}

		// Можно вернуть обновлённый объект, но для простоты — 200 OK
		respondJSON(w, http.StatusOK, map[string]string{"message": "Car updated"})
	}
}

func DeleteCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/cars/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid car ID")
			return
		}

		// Добавь метод Delete в repository
		if !repo.Delete(id) {
			respondError(w, http.StatusNotFound, "Car not found")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "Car deleted"})
	}
}

// ────────────────────────────────────────────────
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
