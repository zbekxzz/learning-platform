package main

import (
	"log"
	"os"

	"platform/internal/auth"
	"platform/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"platform/internal/courses"
	"platform/internal/enrollments"
	"platform/internal/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	database.Connect()

	r := gin.Default()

	api := r.Group("/api")
	auth.RegisterRoutes(api)
	courses.RegisterPublicRoutes(api)

	protected := api.Group("/")
	protected.Use(middleware.JWT())

	courses.RegisterProtectedRoutes(protected)
	enrollments.RegisterRoutes(protected)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
