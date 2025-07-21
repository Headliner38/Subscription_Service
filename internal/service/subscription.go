package service

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Headliner38/Subscription_Service/internal/model"
	"github.com/Headliner38/Subscription_Service/internal/repository"
	"github.com/Headliner38/Subscription_Service/internal/utils"
)

type SubscriptionService struct {
	DB *sql.DB
}

func (s *SubscriptionService) CreateSubscription(
	serviceName string,
	price int,
	userID string,
	startDateStr string,
	endDateStr *string,
) (*model.Subscription, error) {
	log.Printf("[SERVICE] Creating subscription for user: %s, service: %s, price: %d", userID, serviceName, price)

	// Валидация данных
	if serviceName == "" {
		log.Printf("[ERROR] Service name is required")
		return nil, errors.New("service name is required")
	}
	if price <= 0 {
		log.Printf("[ERROR] Price must be positive, got: %d", price)
		return nil, errors.New("price must be positive")
	}
	if userID == "" {
		log.Printf("[ERROR] User ID is required")
		return nil, errors.New("user_id is required")
	}

	// Преобразование дат
	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		log.Printf("[ERROR] Invalid start_date format: %s", startDateStr)
		return nil, errors.New("invalid start_date format, expected MM-YYYY")
	}

	var endDate *time.Time
	if endDateStr != nil && *endDateStr != "" {
		t, err := time.Parse("01-2006", *endDateStr)
		if err != nil {
			log.Printf("[ERROR] Invalid end_date format: %s", *endDateStr)
			return nil, errors.New("invalid end_date format, expected MM-YYYY")
		}
		endDate = &t
		if endDate.Before(startDate) {
			log.Printf("[ERROR] End date cannot be before start date: %s < %s", *endDateStr, startDateStr)
			return nil, errors.New("end_date cannot be before start_date")
		}
	}

	// Генерация UUID
	id := utils.GenerateUUID()
	log.Printf("[SERVICE] Generated UUID: %s", id)

	// Создание структуры подписки
	sub := model.NewSubscription(id, serviceName, price, userID, startDate, endDate)

	// Вызов репозитория для сохранения в БД
	err = repository.CreateSubscription(s.DB, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	if err != nil {
		log.Printf("[ERROR] Failed to save subscription to DB: %v", err)
		return nil, err
	}

	log.Printf("[SUCCESS] Subscription created successfully with ID: %s", id)
	return sub, nil
}

func (s *SubscriptionService) GetSubscription(id string) (*model.Subscription, error) {
	log.Printf("[SERVICE] Getting subscription with ID: %s", id)

	if id == "" {
		log.Printf("[ERROR] ID is required")
		return nil, errors.New("id is required")
	}

	sub, err := repository.GetSubscription(s.DB, id)
	if err != nil {
		log.Printf("[ERROR] Failed to get subscription from DB: %v", err)
		return nil, err
	}

	log.Printf("[SUCCESS] Retrieved subscription with ID: %s", id)
	return sub, nil
}

func (s *SubscriptionService) UpdateSubscription(
	id string,
	serviceName string,
	price int,
	userID string,
	startDateStr string,
	endDateStr *string,
) (*model.Subscription, error) {
	log.Printf("[SERVICE] Updating subscription with ID: %s", id)

	// Валидация
	if id == "" {
		log.Printf("[ERROR] ID is required for update")
		return nil, errors.New("id is required")
	}
	if serviceName == "" {
		log.Printf("[ERROR] Service name is required for update")
		return nil, errors.New("service name is required")
	}
	if price <= 0 {
		log.Printf("[ERROR] Price must be positive for update, got: %d", price)
		return nil, errors.New("price must be positive")
	}
	if userID == "" {
		log.Printf("[ERROR] User ID is required for update")
		return nil, errors.New("user_id is required")
	}

	// Преобразование дат
	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		log.Printf("[ERROR] Invalid start_date format for update: %s", startDateStr)
		return nil, errors.New("invalid start_date format, expected MM-YYYY")
	}

	var endDate *time.Time
	if endDateStr != nil && *endDateStr != "" {
		t, err := time.Parse("01-2006", *endDateStr)
		if err != nil {
			log.Printf("[ERROR] Invalid end_date format for update: %s", *endDateStr)
			return nil, errors.New("invalid end_date format, expected MM-YYYY")
		}
		endDate = &t
		if endDate.Before(startDate) {
			log.Printf("[ERROR] End date cannot be before start date for update: %s < %s", *endDateStr, startDateStr)
			return nil, errors.New("end_date cannot be before start_date")
		}
	}

	// Обновление в БД
	err = repository.UpdateSubscription(s.DB, id, serviceName, price, userID, startDate, endDate)
	if err != nil {
		log.Printf("[ERROR] Failed to update subscription in DB: %v", err)
		return nil, err
	}

	// Возвращаем обновлённую подписку
	sub, err := s.GetSubscription(id)
	if err != nil {
		log.Printf("[ERROR] Failed to get updated subscription: %v", err)
		return nil, err
	}

	log.Printf("[SUCCESS] Updated subscription with ID: %s", id)
	return sub, nil
}

func (s *SubscriptionService) DeleteSubscription(id string) error {
	log.Printf("[SERVICE] Deleting subscription with ID: %s", id)

	if id == "" {
		log.Printf("[ERROR] ID is required for deletion")
		return errors.New("id is required")
	}

	err := repository.DeleteSubscription(s.DB, id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete subscription from DB: %v", err)
		return err
	}

	log.Printf("[SUCCESS] Deleted subscription with ID: %s", id)
	return nil
}

func (s *SubscriptionService) ListSubscriptions() ([]model.Subscription, error) {
	log.Printf("[SERVICE] Listing all subscriptions")

	subscriptions, err := repository.ListSubscriptions(s.DB)
	if err != nil {
		log.Printf("[ERROR] Failed to list subscriptions from DB: %v", err)
		return nil, err
	}

	log.Printf("[SUCCESS] Retrieved %d subscriptions", len(subscriptions))
	return subscriptions, nil
}

func (s *SubscriptionService) CalculateTotalCost(userID, serviceName, startDateStr, endDateStr string) (int, error) {
	log.Printf("[SERVICE] Calculating total cost for user: %s, service: %s, period: %s - %s", userID, serviceName, startDateStr, endDateStr)

	// Преобразование дат
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("01-2006", startDateStr)
		if err != nil {
			log.Printf("[ERROR] Invalid start_date format for total cost: %s", startDateStr)
			return 0, errors.New("invalid start_date format, expected MM-YYYY")
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("01-2006", endDateStr)
		if err != nil {
			log.Printf("[ERROR] Invalid end_date format for total cost: %s", endDateStr)
			return 0, errors.New("invalid end_date format, expected MM-YYYY")
		}
	}

	// Проверка логики дат
	if startDateStr != "" && endDateStr != "" && endDate.Before(startDate) {
		log.Printf("[ERROR] End date cannot be before start date for total cost: %s < %s", endDateStr, startDateStr)
		return 0, errors.New("end_date cannot be before start_date")
	}

	// Вызов репозитория для подсчёта
	totalCost, err := repository.CalculateTotalCost(s.DB, userID, serviceName, startDate, endDate)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate total cost in DB: %v", err)
		return 0, err
	}

	log.Printf("[SUCCESS] Calculated total cost: %d", totalCost)
	return totalCost, nil
}
