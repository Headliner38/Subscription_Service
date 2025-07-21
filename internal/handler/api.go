package handler

import (
	"github.com/Headliner38/Subscription_Service/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, subscriptionService *service.SubscriptionService) {
	subscriptionHandler := &SubscriptionHandler{Service: subscriptionService}

	// Группа маршрутов для подписок
	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.POST("/", subscriptionHandler.CreateSubscription)
		subscriptions.GET("/", subscriptionHandler.ListSubscriptions)
		subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
		subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
		subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
		subscriptions.GET("/total", subscriptionHandler.CalculateTotalCost)
	}
}
