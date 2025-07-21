package handler

import (
	"log"
	"net/http"

	"github.com/Headliner38/Subscription_Service/internal/service"
	"github.com/gin-gonic/gin"
)

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Netflix" binding:"required"`
	Price       int    `json:"price" example:"999" binding:"required,gt=0"`
	UserID      string `json:"user_id" example:"user123" binding:"required"`
	StartDate   string `json:"start_date" example:"01-2024" binding:"required"`
	EndDate     string `json:"end_date,omitempty" example:"12-2024"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Netflix" binding:"required"`
	Price       int    `json:"price" example:"999" binding:"required,gt=0"`
	UserID      string `json:"user_id" example:"user123" binding:"required"`
	StartDate   string `json:"start_date" example:"01-2024" binding:"required"`
	EndDate     string `json:"end_date,omitempty" example:"12-2024"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

type TotalCostResponse struct {
	TotalCost   int    `json:"total_cost" example:"2997"`
	UserID      string `json:"user_id" example:"user123"`
	ServiceName string `json:"service_name" example:"Netflix"`
	StartDate   string `json:"start_date" example:"01-2024"`
	EndDate     string `json:"end_date" example:"12-2024"`
}

type SubscriptionHandler struct {
	Service *service.SubscriptionService
}

// CreateSubscription godoc
// @Summary Создать подписку
// @Description Создаёт новую подписку для пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	log.Printf("[HANDLER] Creating subscription")

	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var endDate *string
	if req.EndDate != "" {
		endDate = &req.EndDate
	}

	sub, err := h.Service.CreateSubscription(req.ServiceName, req.Price, req.UserID, req.StartDate, endDate)
	if err != nil {
		log.Printf("[ERROR] Failed to create subscription: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Created subscription with ID: %s", sub.ID)
	c.JSON(http.StatusCreated, sub)
}

// GetSubscription godoc
// @Summary Получить подписку
// @Description Получает подписку по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[HANDLER] Getting subscription with ID: %s", id)

	sub, err := h.Service.GetSubscription(id)
	if err != nil {
		log.Printf("[ERROR] Failed to get subscription %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Retrieved subscription with ID: %s", id)
	c.JSON(http.StatusOK, sub)
}

// ListSubscriptions godoc
// @Summary Список подписок
// @Description Получает список всех подписок
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	log.Printf("[HANDLER] Listing all subscriptions")

	subscriptions, err := h.Service.ListSubscriptions()
	if err != nil {
		log.Printf("[ERROR] Failed to list subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Retrieved %d subscriptions", len(subscriptions))
	c.JSON(http.StatusOK, subscriptions)
}

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Description Обновляет существующую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body UpdateSubscriptionRequest true "Новые данные подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[HANDLER] Updating subscription with ID: %s", id)

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Invalid request body for update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var endDate *string
	if req.EndDate != "" {
		endDate = &req.EndDate
	}

	sub, err := h.Service.UpdateSubscription(id, req.ServiceName, req.Price, req.UserID, req.StartDate, endDate)
	if err != nil {
		log.Printf("[ERROR] Failed to update subscription %s: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Updated subscription with ID: %s", id)
	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[HANDLER] Deleting subscription with ID: %s", id)

	err := h.Service.DeleteSubscription(id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete subscription %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Deleted subscription with ID: %s", id)
	c.Status(http.StatusNoContent)
}

// CalculateTotalCost godoc
// @Summary Подсчитать общую стоимость
// @Description Подсчитывает общую стоимость подписок с фильтрацией
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string false "Начальная дата (MM-YYYY)"
// @Param end_date query string false "Конечная дата (MM-YYYY)"
// @Success 200 {object} TotalCostResponse
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	// Получаем параметры из query string
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	log.Printf("[HANDLER] Calculating total cost for user: %s, service: %s, period: %s - %s", userID, serviceName, startDate, endDate)

	// Вызываем сервис для подсчёта
	totalCost, err := h.Service.CalculateTotalCost(userID, serviceName, startDate, endDate)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate total cost: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Calculated total cost: %d", totalCost)
	// Возвращаем результат
	c.JSON(http.StatusOK, gin.H{
		"total_cost":   totalCost,
		"user_id":      userID,
		"service_name": serviceName,
		"start_date":   startDate,
		"end_date":     endDate,
	})
}
