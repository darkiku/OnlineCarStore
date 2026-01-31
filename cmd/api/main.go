package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/teamserik/online-car-store/internal/handler"
	"github.com/teamserik/online-car-store/internal/repository"
)

var (
	carRepo   = repository.NewInMemoryCarRepository()
	viewCount int64 // для демонстрации горутины
)

func main() {
	// Фоновая горутина — требование задания (concurrency)
	go backgroundStatsPrinter()

	mux := http.NewServeMux()

	// Роуты (все JSON)
	mux.HandleFunc("POST /api/cars", handler.CreateCar(carRepo))
	mux.HandleFunc("GET /api/cars", handler.ListCars(carRepo))
	mux.HandleFunc("GET /api/cars/", handler.GetCar(carRepo))    // /api/cars/5
	mux.HandleFunc("PUT /api/cars/", handler.UpdateCar(carRepo)) // /api/cars/5
	mux.HandleFunc("DELETE /api/cars/", handler.DeleteCar(carRepo))

	fmt.Println("Car Store API running on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))
}

// Пример background worker (горутина + atomic)
func backgroundStatsPrinter() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		count := atomic.LoadInt64(&viewCount)
		fmt.Printf("[BG] Total car detail views so far: %d\n", count)
	}
}
