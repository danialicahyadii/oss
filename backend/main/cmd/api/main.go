package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"oss.kftd.co.id/v2/main/routes"
	"oss.kftd.co.id/v2/main/shared/db"
	"oss.kftd.co.id/v2/user-management/seeder"
)

func main() {

	if err := godotenv.Load("../infrastructure/.env"); err != nil {
		log.Println("warning: .env not found")
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatal("database connection failed:", err)
	}

	log.Println("database connected")

	// Run seeder
	if err := seeder.SeedUserManagement(database); err != nil {
		log.Fatal("user management seeder failed:", err)
	}

	log.Println("user management seeded")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	routes.Setup(router, database)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("server running on port", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
