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

	empStore, err := repositories.NewEmp(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer empStore.DisconnectEmp(timeoutContext)

	compStore, err := repositories.NewCompanyRepo(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer compStore.DisconnectComp(timeoutContext)

	userStore, err := repositories.NewUserRepo(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer userStore.DisconnectUser(timeoutContext)

	empStore.PingEmp()
	userStore.PingUser()
	compStore.PingComp()

	employmentHandler := handlers.NewEmploymentHandler(logger, empStore)
	jobHandler := handlers.NewJobHandler(logger, empStore)
	newsHandler := handlers.NewNewsHandler(logger, empStore)
	jobAdHandler := handlers.NewJobAdHandler(logger, empStore)
	reviewOfCompanyHandler := handlers.NewReviewOfCompanyHandler(logger, empStore)
	applicantHandler := handlers.NewApplicantHandler(logger, empStore)
	companyHandler := handlers.NewCompanyHandler(logger, compStore)
	cvHandler := handlers.NewCvHandler(logger, compStore)

	router := gin.New()
	router.Use(jobAdHandler.JobadCORSMiddleware())
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

	router.GET("/reviewofcompany/company/{company_id}", reviewOfCompanyHandler.GetReviewsOfCompany)
	router.GET("/reviewofcompany/{review_id}", reviewOfCompanyHandler.GetReviewById)
	router.POST("/reviewofcompany", reviewOfCompanyHandler.PostReview)
	router.PUT("/reviewofcompany/{review_id}", reviewOfCompanyHandler.EditReview)
	router.DELETE("/reviewofcompany/{review_id}", reviewOfCompanyHandler.DeleteReviewById)

	router.GET("/applicant", applicantHandler.GetApplicants)
	router.GET("/applicant/company/{ucn}", applicantHandler.GetApplicantsForCompany)
	router.GET("/applicant/{ucn}", applicantHandler.GetApplicantByUcn)
	router.POST("/applicant/{ucn}", applicantHandler.PostApplicant)
	router.DELETE("/applicant/{applicant_id}", applicantHandler.DeleteApplicantById)

	router.POST("/company", companyHandler.CreateCompany)
	router.GET("/company", companyHandler.GetCompanies)
	router.GET("/company/{id}", companyHandler.FindCompanyById)
	router.GET("/company/owner/{owner}", companyHandler.GetCompanyByOwnerUcn)

	router.GET("/resume/{ucn}", cvHandler.FindCvByUcn)
	router.POST("/resume/{ucn}", cvHandler.PostCv)

	router.POST("/employee", employmentHandler.EmployApplicant)

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
