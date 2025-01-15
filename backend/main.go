package main

import (
	"context"
	"employment-service/handlers"
	"employment-service/repositories"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	jobHandler := handlers.NewJobHandler(logger, store)

	router := gin.New()
	router.Use(employmentHandler.CORSMiddleware())
	router.GET("/jobs", jobHandler.GetJobs)
	router.GET("/job/{job_id}", jobHandler.GetJobById)
	router.POST("/job", jobHandler.PostJob)
	router.PUT("/job/{job_id}", jobHandler.EditJob)
	router.DELETE("/job/{job_id}", jobHandler.DeleteJobById)

	router.Run(":" + port)

	server := http.Server{
		Addr:         ":" + port,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)
}
