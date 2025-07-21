// @title Subscription Service API
// @version 1.0
// @description REST API для управления подписками пользователей
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"database/sql"
	"log"

	_ "github.com/Headliner38/Subscription_Service/docs" // Swagger docs
	"github.com/Headliner38/Subscription_Service/internal/config"
	"github.com/Headliner38/Subscription_Service/internal/handler"
	"github.com/Headliner38/Subscription_Service/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	log.Printf("[MAIN] Starting Subscription Service...")

	// Загрузка конфигурацию
	cfg := config.LoadConfig()
	log.Printf("[MAIN] Configuration loaded successfully")

	// Подключение к БД
	connStr := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " sslmode=disable"
	log.Printf("[MAIN] Connecting to database at %s:%s", cfg.DBHost, cfg.DBPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("[FATAL] Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Проверка подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("[FATAL] Failed to ping database: %v", err)
	}
	log.Printf("[MAIN] Database connection established successfully")

	// Инициализируем сервисы и обработчики
	subscriptionService := &service.SubscriptionService{DB: db}
	log.Printf("[MAIN] Services initialized")

	// роутер
	r := gin.New() // gin.New() для кастомного логирования

	// middleware для логирования
	r.Use(handler.LoggerMiddleware())
	r.Use(gin.Recovery()) // recovery middleware

	handler.SetupRoutes(r, subscriptionService)
	log.Printf("[MAIN] Routes configured")

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Printf("[MAIN] Swagger UI available at /swagger/index.html")

	// Запуск сервера
	log.Printf("[MAIN] Server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("[FATAL] Failed to start server: %v", err)
	}
}
