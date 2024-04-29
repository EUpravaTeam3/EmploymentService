package main

import (
	"context"
	"employment-service/handlers"
	"employment-service/repositories"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[employment-store] ", log.LstdFlags)

	store, err := repositories.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	store.Ping()

	employmentHandler := handlers.NewEmploymentHandler(logger, store)

	router := gin.New()
	router.Use(employmentHandler.CORSMiddleware())
	router.Use(employmentHandler.GenerateDiploma)
}
