package repository

import (
	"sync"
	"time"

	"github.com/teamserik/online-car-store/internal/model"
)

type CarRepository interface {
	Create(car *model.Car)
	GetByID(id int) (*model.Car, bool)
	List() []*model.Car
	Update(id int, input model.UpdateCarInput) bool
	Delete(id int) bool
	// Delete можно добавить позже, если останется время
}

type inMemoryCarRepository struct {
	cars   map[int]*model.Car
	mu     sync.RWMutex
	nextID int
}

func NewInMemoryCarRepository() CarRepository {
	return &inMemoryCarRepository{
		cars:   make(map[int]*model.Car),
		nextID: 1,
	}
}

func (r *inMemoryCarRepository) Create(car *model.Car) {
	r.mu.Lock()
	defer r.mu.Unlock()

	car.ID = r.nextID
	r.nextID++
	now := time.Now()
	car.CreatedAt = now
	car.UpdatedAt = now

	r.cars[car.ID] = car
}

func (r *inMemoryCarRepository) GetByID(id int) (*model.Car, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	car, exists := r.cars[id]
	return car, exists
}

func (r *inMemoryCarRepository) List() []*model.Car {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Car, 0, len(r.cars))
	for _, car := range r.cars {
		result = append(result, car)
	}
	return result
}

func (r *inMemoryCarRepository) Update(id int, input model.UpdateCarInput) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	car, exists := r.cars[id]
	if !exists {
		return false
	}

	if input.Price != nil {
		car.Price = *input.Price
	}
	if input.Mileage != nil {
		car.Mileage = *input.Mileage
	}
	if input.Description != nil {
		car.Description = *input.Description
	}

	car.UpdatedAt = time.Now()
	return true
}

func (r *inMemoryCarRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cars[id]; !exists {
		return false
	}
	delete(r.cars, id)
	return true
}
