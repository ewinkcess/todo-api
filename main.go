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
	database.DB.AutoMigrate(&domain.User{}, &domain.Todo{}, &domain.Category{})

	todoRepo := repository.NewTodoRepository(database.DB)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	categoryRepo := repository.NewCategoryRepository(database.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	router := gin.Default()
	router.Use(middleware.LoggerMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandler())

	auth := router.Group("/auth")
	auth.Use(middleware.RateLimiter(5, time.Minute))
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	todos := router.Group("/todos")
	todos.Use(middleware.AuthMiddleware())
	{
		todos.GET("", todoHandler.GetAllTodos)
		todos.GET("/:id", todoHandler.GetTodoByID)
		todos.POST("", todoHandler.CreateTodo)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/me", userHandler.GetProfile)
		users.PUT("/me", userHandler.UpdateProfile)
		users.PUT("/me/password", userHandler.UpdatePassword)
	}

	categories := router.Group("/categories")
	categories.Use(middleware.AuthMiddleware())
	{
		categories.GET("", categoryHandler.GetAllCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		categories.POST("", categoryHandler.CreateCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	log.Println("Restoran Todo API buka di port", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Restoran gagal dibuka:", err)
	}
}
