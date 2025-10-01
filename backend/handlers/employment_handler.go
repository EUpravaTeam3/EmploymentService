package handlers

import (
	"bytes"
	"context"
	"employment-service/domain"
	"employment-service/repositories"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmploymentHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewEmploymentHandler(l *log.Logger, r *repositories.EmploymentRepo) *EmploymentHandler {
	return &EmploymentHandler{l, r}
}

type AcceptedApplicant struct {
	PositionName   string        `bson:"position_name" json:"position_name"`
	AdTitle        string        `bson:"ad_title" json:"ad_title"`
	JobAdId        string        `bson:"job_ad_id" json:"job_ad_id"`
	CitizenUcn     string        `bson:"citizen_ucn" json:"citizen_ucn"`
	Name           string        `bson:"name" json:"name"`
	Email          string        `bson:"email" json:"email"`
	Education      []interface{} `bson:"education" json:"education"`
	WorkExperience []interface{} `bson:"work_experience" json:"work_experience"`
	Description    string        `bson:"description" json:"description"`
}

func (e *EmploymentHandler) EmployApplicant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cookie, err := c.Cookie("SESSION_ID")

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	var applicant AcceptedApplicant
	var user domain.User
	var job domain.Job
	employee := &domain.Employee{}

	jobHandler := JobHandler{repo: e.repo}

	if err := c.ShouldBindJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("ACCEPT EMPLOYEE")
	fmt.Println(applicant.JobAdId)

	jobAdId, err := primitive.ObjectIDFromHex(applicant.JobAdId)

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	job, err = jobHandler.GetJobByJobAdId(jobAdId)

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	employee.CitizenUCN = user.UCN
	employee.JobId = job.Id
	employee.StartDate = time.Now()

	err = e.repo.AddEmployeeToCompany(*employee) // add new employee to employees database

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	var cvColl = e.repo.GetCollection("employmentdb", "cvs")
	var companyColl = e.repo.GetCollection("employmentdb", "companies")
	var jobAdsColl = e.repo.GetCollection("employmentdb", "jobads")
	var applicantsColl = e.repo.GetCollection("employmentdb", "applicants")

	var company domain.Company

	err = companyColl.FindOne(ctx, bson.M{"_id": job.CompanyId}).Decode(&company)

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	/*payloadApr := AprPayload{
		RegistrationNumber: company.IdNumber,
		EmployeeUcn:        applicant.CitizenUcn,
	}*/

	ssoPayload := SsoPayload{
		Role:      "employee",
		Initiator: company.OwnerUcn,
		SessionId: cookie,
	}

	// add employe to APR company list
	// if err := sendToService(c, payloadApr, "http://aprcroso_service:8005/api/companies", false); err != nil {
	// 	e.logger.Println(err)
	// 	c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	// 	return
	// }

	// change users role from citizen to employee in SSO
	fmt.Println(applicant.CitizenUcn)
	if err := sendToService(c, ssoPayload, "http://sso_service:9090/user/employment/"+applicant.CitizenUcn, true); err != nil {
		e.logger.Println(err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var workExperience string = job.PoistionName + " (" + time.Now().Month().String() +
		" " + strconv.Itoa(time.Now().Year()) + " - "

	_, err = cvColl.UpdateOne(ctx, bson.M{"citizen_ucn": applicant.CitizenUcn},
		bson.M{"$push": bson.M{"work_experience": workExperience}}) // add new job start date to cv

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	_, err = applicantsColl.DeleteMany(ctx, bson.D{{Key: "job_ad_id", Value: jobAdId}}) // delete other pending applicants for this job

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	_, err = jobAdsColl.DeleteOne(ctx, bson.D{{Key: "_id", Value: jobAdId}}) // delete the job ad for which this user applied to

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}
}

type AprPayload struct {
	RegistrationNumber string `json:"registrationNumber"`
	EmployeeUcn        string `json:"employeeUcn"`
}

type SsoPayload struct {
	Role      string `bson:"role" json:"role"`
	Initiator string `bson:"initiator" json:"initiator"`
	SessionId string `bson:"session" json:"session"`
}

func sendToService(c *gin.Context, payload interface{}, url string, isCookie bool) error {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Errorf("failed to encode body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Errorf("failed to create request: %w", err)
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
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("remote service returned %d: %s", resp.StatusCode, string(body))
	}

	return nil

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
