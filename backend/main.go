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
	newsHandler := handlers.NewNewsHandler(logger, store)
	jobAdHandler := handlers.NewJobAdHandler(logger, store)

	router := gin.New()
	router.Use(employmentHandler.CORSMiddleware())
	router.GET("/jobs", jobHandler.GetJobs)
	router.GET("/job/{job_id}", jobHandler.GetJobById)
	router.POST("/job", jobHandler.PostJob)
	router.PUT("/job/{job_id}", jobHandler.EditJob)
	router.DELETE("/job/{job_id}", jobHandler.DeleteJobById)

	router.GET("/news", newsHandler.GetNews)
	router.GET("/news/{news_id}", newsHandler.GetNewsById)
	router.POST("/news", newsHandler.PostNews)
	router.PUT("/news/{news_id}", newsHandler.EditNews)
	router.DELETE("/news/{news_id}", newsHandler.DeleteNewsById)

	router.GET("/jobad", jobAdHandler.GetJobAds)
	router.GET("/jobad/{jobad_id}", jobAdHandler.GetJobAdById)
	router.POST("/jobad", jobAdHandler.PostJobAd)
	router.PUT("/jobad/{jobad_id}", jobAdHandler.EditJobAd)
	router.DELETE("/jobad/{jobad_id}", jobAdHandler.DeleteJobAdById)

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
