package handlers

import (
	"employment-service/domain"
	"employment-service/repositories"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type EmploymentHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewEmploymentHandler(l *log.Logger, r *repositories.EmploymentRepo) *EmploymentHandler {
	return &EmploymentHandler{l, r}
}

func (e *EmploymentHandler) EmployApplicant(c *gin.Context) {

	var applicant *domain.Applicant
	var user domain.User
	var job domain.Job
	var employee *domain.Employee

	var userHandler UserHandler
	var jobHandler JobHandler

	if err := c.ShouldBindJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ = userHandler.FindUserByUcn("2309998800186") // ***for testing purposes only, change later to CV id!!

	job, _ = jobHandler.GetJobByJobAdId(applicant.JobAdId)

	employee.CitizenUCN = user.UCN
	employee.JobId = job.Id
	employee.StartDate = time.Now()

	e.repo.AddEmployeeToCompany(*employee)
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
