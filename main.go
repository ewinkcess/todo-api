// ============================================
// FILE: main.go
// ============================================

// @title           Todo API
// @version         1.0
// @description     REST API untuk manajemen Todo dengan autentikasi JWT

// @contact.name    Ewink
// @contact.email   ewink@gmail.com

// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Masukkan token dengan format: Bearer <token>

package main

import (
	"log"
	"time"

	"todo-api/config"
	"todo-api/database"

	_ "todo-api/docs"
	"todo-api/internal/domain"
	"todo-api/internal/handler"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	middleware.InitLogger()
	defer middleware.Logger.Sync()
	cfg := config.LoadConfig()
	database.ConnectDatabase(cfg)
	database.DB.AutoMigrate(&domain.User{})

	todoRepo := repository.NewTodoRepository(database.DB)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	router.Use(middleware.LoggerMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	auth.Use(middleware.RateLimiter(5, time.Minute))
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Langkah 8 — Rute protected (wajib lewat middleware)
	todos := router.Group("/todos")
	todos.Use(middleware.AuthMiddleware())
	{
		todos.GET("", todoHandler.GetAllTodos)
		todos.GET("/:id", todoHandler.GetTodoByID)
		todos.POST("", todoHandler.CreateTodo)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}

	// Langkah 9 — Buka restoran
	log.Println("Restoran Todo API buka di port", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Restoran gagal dibuka:", err)
	}
}
