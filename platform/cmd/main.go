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
	"platform/internal/modules"
	"platform/internal/tests"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	database.Connect()

	r := gin.Default()
	r.RedirectTrailingSlash = false

	api := r.Group("/api/v1")
	auth.RegisterRoutes(api)
	courses.RegisterPublicRoutes(api)

	protected := api.Group("/")
	protected.Use(middleware.JWT())

	courses.RegisterProtectedRoutes(protected)
	enrollments.RegisterRoutes(protected)
	modules.RegisterRoutes(protected)
	tests.RegisterRoutes(protected)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
