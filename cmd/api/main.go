package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/omniflare/go_starter/internals/config"
	"github.com/omniflare/go_starter/internals/db"
	"github.com/omniflare/go_starter/internals/http/handlers"

	// middleware "github.com/omniflare/go_starter/internals/http/middlewares"
	repositories "github.com/omniflare/go_starter/internals/http/repositories/auth"
	service "github.com/omniflare/go_starter/internals/services/auth"
)

func main() {
	envConfig := config.NewEnvConfig()

	db := db.Init(envConfig, db.AutoMigrate)
	app := gin.Default()
	authRepository := repositories.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	api := app.Group("/api")
	auth := api.Group("/auth")
	handlers.NewAuthHandler(auth, *authService)

	{
		// auth.POST("/login", authService.Login)
		// auth.POST("/register", authService.Register)

		// protected := auth.Group("")
		// protected.Use(middleware.AuthMiddleware())
		// {
		//     protected.POST("/verify-email", authService.VerifyEmail)
		//     protected.POST("/forgot-password", authService.ForgotPassword)
		//     protected.POST("/reset-password", authService.ResetPassword)
		// }
	}
	port := envConfig.PORT
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
