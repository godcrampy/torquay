package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mode := os.Getenv("ENV")
	gin.SetMode(mode)

	r := gin.Default()

	token := 1

	r.GET("/api/v1/token", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})

		token += 1
	})

	port := os.Getenv("PORT")
	log.Printf("INFO: Starting server on port %s\n", port)
	log.Printf("INFO: Server running in %s mode\n", mode)
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
