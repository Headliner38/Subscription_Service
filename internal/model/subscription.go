package model

import (
	"time"
)

// Subscription представляет подписку пользователя
// @Description Модель подписки пользователя
type Subscription struct {
	ID          string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" db:"id"`
	ServiceName string     `json:"service_name" example:"Netflix" db:"service_name"`
	Price       int        `json:"price" example:"999" db:"price"`
	UserID      string     `json:"user_id" example:"user123" db:"user_id"`
	StartDate   time.Time  `json:"start_date" example:"2024-01-01T00:00:00Z" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"2024-12-31T00:00:00Z" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z" db:"updated_at"`
}

// NewSubscription создаёт новую подписку
func NewSubscription(id, serviceName string, price int, userID string, startDate time.Time, endDate *time.Time) *Subscription {
	now := time.Now()
	return &Subscription{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
