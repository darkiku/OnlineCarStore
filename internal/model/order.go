package model

import "time"

type Order struct {
	ID        int       `json:"id"`
	CarID     int       `json:"car_id"`
	UserID    int       `json:"user_id"` // пока заглушка, auth нет
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"` // pending, confirmed, cancelled
}
