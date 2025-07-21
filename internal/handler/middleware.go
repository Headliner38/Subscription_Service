package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware логирует все HTTP-запросы
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Логируем запрос с временем выполнения
		log.Printf("[HTTP] %s | %s | %d | %v | %s | %s",
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.ErrorMessage,
		)
		return ""
	})
}
