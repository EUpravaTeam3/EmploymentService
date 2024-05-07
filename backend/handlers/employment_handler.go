package handlers

import (
	"employment-service/domain"
	"employment-service/repositories"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type EmploymentHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func NewEmploymentHandler(l *log.Logger, r *repositories.EmploymentRepo) *EmploymentHandler {
	return &EmploymentHandler{l, r}
}

func (e *EmploymentHandler) GenerateDiploma(c *gin.Context) {

	var diploma *domain.Diploma

	if err := c.BindJSON(&diploma); err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := e.repo.AddDiplomaToCV(diploma); err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diploma)

}

func (e *EmploymentHandler) PostJobAd(c *gin.Context) {

	var jobAd *domain.JobAd

	if err := c.BindJSON(&jobAd); err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var company, err = e.repo.FindCompanyByJobId(jobAd.JobId)

	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	companyURL := "http://localhost:..." // path to kolega's controller

	req, err := http.NewRequest(http.MethodGet, companyURL+string(company.IdNumber), nil)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while processing your request"})
		return
	}

	res, errClient := http.DefaultClient.Do(req)

	if errClient != nil || res.StatusCode != 200 {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "APR refused to check the company status, try again later"})
		return
	}

	defer res.Body.Close()

	var status string

	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Unexpected error"})
		return
	}

	if status != "ACTIVE" {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Your company is not active!"})
	}

	if err := e.repo.InsertJobAd(jobAd); err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Unexpected error"})
	}

	c.JSON(http.StatusOK, nil)
}

func (e *EmploymentHandler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
