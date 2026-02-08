package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/teamserik/online-car-store/internal/model"
	"github.com/teamserik/online-car-store/internal/repository"
)

func CreateCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.CreateCarInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		car := &model.Car{
			Make:         input.Make,
			Model:        input.Model,
			Year:         input.Year,
			Price:        input.Price,
			Mileage:      input.Mileage,
			BodyType:     input.BodyType,
			FuelType:     input.FuelType,
			Transmission: input.Transmission,
			Color:        input.Color,
			HorsePower:   input.HorsePower,
			EngineSize:   input.EngineSize,
			Description:  input.Description,
			ImageURL:     input.ImageURL,
		}

		ctx := context.Background()
		if err := repo.Create(ctx, car); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(car)
	}
}

func ListCars(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filter model.FilterParams

		// Парсинг query параметров для фильтрации
		query := r.URL.Query()

		if make := query.Get("make"); make != "" {
			filter.Make = &make
		}
		if bodyType := query.Get("body_type"); bodyType != "" {
			filter.BodyType = &bodyType
		}
		if fuelType := query.Get("fuel_type"); fuelType != "" {
			filter.FuelType = &fuelType
		}
		if transmission := query.Get("transmission"); transmission != "" {
			filter.Transmission = &transmission
		}

		ctx := context.Background()
		cars, err := repo.List(ctx, &filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cars)
	}
}

func GetCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/cars/")

		ctx := context.Background()
		car, err := repo.GetByID(ctx, id)
		if err != nil {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(car)
	}
}

func UpdateCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/cars/")

		var input model.UpdateCarInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		if err := repo.Update(ctx, id, input); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Car updated successfully"})
	}
}

func DeleteCar(repo repository.CarRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/cars/")

		ctx := context.Background()
		if err := repo.Delete(ctx, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Car deleted successfully"})
	}
}
