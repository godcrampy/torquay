package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/godcrampy/torquay/pkg/counter"
	"github.com/godcrampy/torquay/pkg/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// time.Sleep(time.Second * 10)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mode := os.Getenv("ENV")
	gin.SetMode(mode)

	servers := []string{os.Getenv("ZOOKEEPER_SERVERS")}
	zkPath := "/counter"

	c, err := counter.NewCounter(servers, zkPath)
	if err != nil {
		log.Fatalf("Unable to connect to ZooKeeper: %v", err)
	}
	defer c.Close()

	h := handlers.NewHandler(c)

	r := gin.Default()
	r.GET("/api/v1/token", h.GetToken)

	port := os.Getenv("PORT")
	log.Printf("INFO: Starting server on port %s\n", port)
	log.Printf("INFO: Server running in %s mode\n", mode)
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
