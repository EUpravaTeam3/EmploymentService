package handlers

import (
	"bytes"
	"context"
	"employment-service/domain"
	"employment-service/repositories"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type EmploymentHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewEmploymentHandler(l *log.Logger, r *repositories.EmploymentRepo) *EmploymentHandler {
	return &EmploymentHandler{l, r}
}

func (e *EmploymentHandler) EmployApplicant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var applicant *domain.Applicant
	var user domain.User
	var job domain.Job
	var employee *domain.Employee

	var jobHandler JobHandler

	if err := c.ShouldBindJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, _ = jobHandler.GetJobByJobAdId(applicant.JobAdId)

	employee.CitizenUCN = user.UCN
	employee.JobId = job.Id
	employee.StartDate = time.Now()

	err := e.repo.AddEmployeeToCompany(*employee) // add new employee to employees database

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}

	var cvColl = e.repo.GetCollection("employmentdb", "cvs")
	var companyColl = e.repo.GetCollection("employmentdb", "companies")
	var jobAdsColl = e.repo.GetCollection("employmentdb", "jobads")
	var applicantsColl = e.repo.GetCollection("employmentdb", "applicants")

	var cv domain.CV
	var company domain.Company

	err = companyColl.FindOne(ctx, bson.M{"_id": job.CompanyId}).Decode(&company)

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}

	err = cvColl.FindOne(ctx, bson.M{"_id": applicant.CVId}).Decode(&cv)

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}

	payloadApr := AprPayload{
		RegistrationNumber: company.IdNumber,
		EmployeeUcn:        cv.CitizenUCN,
	}

	// change the URL from localhost to aprcroso_service for Docker
	sendToService(c, payloadApr, "http://localhost:8005/employees", false)

	// sendToService(c, , "http://sso_service:9090...") // for changing roles

	var workExperience string = job.PoistionName + " " + time.Now().Month().String() +
		"/" + strconv.Itoa(time.Now().Year()) + " - "

	_, err = cvColl.UpdateOne(ctx, bson.M{"_id": applicant.CVId},
		bson.M{"$push": bson.M{"work_experience": workExperience}}) // add new job start date to cv

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}

	_, err = applicantsColl.DeleteMany(ctx, bson.M{"job_ad_id": applicant.JobAdId}) // delete other pending applicants for this job

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}

	_, err = jobAdsColl.DeleteOne(ctx, bson.M{"_id": applicant.JobAdId}) // delete the job ad for which this user applied to

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}
}

type AprPayload struct {
	RegistrationNumber string `json:"registrationNumber"`
	EmployeeUcn        string `json:"employeeUcn"`
}

func sendToService(c *gin.Context, payload interface{}, url string, isCookie bool) {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode body"})
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	if isCookie {
		if cookie, err := c.Cookie("SESSION_ID"); err == nil {
			req.AddCookie(&http.Cookie{
				Name:  "SESSION_ID",
				Value: cookie,
			})
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	c.JSON(resp.StatusCode, gin.H{"message": "Body + cookie forwarded"})

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
