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
	"go.mongodb.org/mongo-driver/mongo"
)

type EmploymentHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewEmploymentHandler(l *log.Logger, r *repositories.EmploymentRepo) *EmploymentHandler {
	return &EmploymentHandler{l, r}
}

type AcceptedApplicant struct {
	PositionName    string        `bson:"position_name" json:"position_name"`
	AdTitle         string        `bson:"ad_title" json:"ad_title"`
	JobAdId         string        `bson:"job_ad_id" json:"job_ad_id"`
	CitizenUcn      string        `bson:"citizen_ucn" json:"citizen_ucn"`
	Name            string        `bson:"name" json:"name"`
	Email           string        `bson:"email" json:"email"`
	CvID            string        `bson:"cv_id" json:"cv_id"`
	CompanyOwnerUcn string        `bson:"company_owner_ucn" json:"company_owner_ucn"`
	Education       []interface{} `bson:"education" json:"education"`
	WorkExperience  []interface{} `bson:"work_experience" json:"work_experience"`
	Description     string        `bson:"description" json:"description"`
}

type EmployeeWithJob struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	CitizenUCN     string             `bson:"citizen_ucn" json:"citizen_ucn"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	EndDate        *time.Time         `bson:"end_date,omitempty" json:"end_date,omitempty"`
	EmployerReview string             `bson:"employer_review,omitempty" json:"employer_review,omitempty"`
	PositionName   string             `bson:"position_name" json:"position_name"`
	Pay            float64            `bson:"pay" json:"pay"`
	CompanyName    string             `bson:"company_name" json:"company_name"`
}

func (e *EmploymentHandler) GetEmployees(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var employees domain.Employees

	empColl := e.repo.GetCollection("employmentdb", "employees")

	empCursor, err := empColl.Find(ctx, bson.M{})
	if err != nil {
		e.logger.Println(err)
		return
	}

	if err = empCursor.All(ctx, &employees); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Fatal(err)
		return
	}

	err = employees.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (e *EmploymentHandler) GetEmployeeByUcn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ucn := c.Param("ucn")

	empColl := e.repo.GetCollection("employmentdb", "employees")

	pipeline := mongo.Pipeline{
		// Match employee(s) by citizen_ucn
		{{"$match", bson.D{{"citizen_ucn", ucn}}}},

		// Lookup the Job for the employee
		{{"$lookup", bson.D{
			{"from", "jobs"},
			{"localField", "job_id"},
			{"foreignField", "_id"},
			{"as", "job"},
		}}},

		// Unwind job (since employee has only one job_id)
		{{"$unwind", "$job"}},

		// Lookup the Company for that job
		{{"$lookup", bson.D{
			{"from", "companies"},
			{"localField", "job.company_id"},
			{"foreignField", "_id"},
			{"as", "company"},
		}}},

		// Unwind company
		{{"$unwind", "$company"}},

		// Project final shape
		{{"$project", bson.D{
			{"_id", 1},
			{"citizen_ucn", 1},
			{"start_date", 1},
			{"end_date", 1},
			{"employer_review", 1},
			{"position_name", "$job.position_name"},
			{"pay", "$job.pay"},
			{"company_name", "$company.name"},
		}}},
	}

	cursor, err := empColl.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}
	var employees []EmployeeWithJob
	if err = cursor.All(ctx, &employees); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
	}
	c.JSON(http.StatusOK, employees)
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

	employee.CitizenUCN = applicant.CitizenUcn
	employee.JobId = job.Id
	employee.StartDate = time.Now()

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	var cvColl = e.repo.GetCollection("employmentdb", "cvs")
	var companyColl = e.repo.GetCollection("employmentdb", "companies")
	var jobAdsColl = e.repo.GetCollection("employmentdb", "jobAds")
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

	payloadApr := AprPayload{
		Name:            applicant.Name,
		Position:        job.PoistionName,
		CreatedByUserId: applicant.CompanyOwnerUcn,
	}

	ssoPayload := SsoPayload{
		Role:      "employee",
		Initiator: company.OwnerUcn,
		SessionId: cookie,
	}

	// change users role from citizen to employee in SSO
	fmt.Println(applicant.CitizenUcn)
	if err := SendToService(c, ssoPayload, "http://sso_service:9090/user/employment/"+applicant.CitizenUcn, true); err != nil {
		e.logger.Println(err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// add employe to APR company list
	if err := SendToService(c, payloadApr, "http://aprcroso:8005/adding-employee/request", false); err != nil {
		e.logger.Println(err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	_, err = e.repo.AddEmployeeToCompany(*employee) // add new employee to employees database

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
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

	objectId, err := primitive.ObjectIDFromHex(applicant.CvID)
	if err != nil {
		log.Println("Invalid id")
	}

	_, err = applicantsColl.DeleteMany(ctx, bson.D{{Key: "cv_id", Value: objectId}}) // delete other applications where this user applied

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

// type AprPayload struct {
// 	RegistrationNumber string `json:"registrationNumber"`
// 	EmployeeUcn        string `json:"employeeUcn"`
// }

type AprPayload struct {
	Name            string `bson:"name" json:"name"`
	Position        string `bson:"position" json:"position"`
	CreatedByUserId string `bson:"createdByUserId" json:"createdByUserId"`
}

type SsoPayload struct {
	Role      string `bson:"role" json:"role"`
	Initiator string `bson:"initiator" json:"initiator"`
	SessionId string `bson:"session" json:"session"`
}

func SendToService(c *gin.Context, payload interface{}, url string, isCookie bool) error {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Errorf("failed to encode body: %w", err)
	}

	fmt.Println("Sending to APR:", url, "payload:", string(jsonData))

	buf := bytes.NewBuffer(jsonData)
	fmt.Println("Buffer contents:", buf.String())

	req, err := http.NewRequest("POST", url, buf)

	if err != nil {
		fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

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

func (e *EmploymentHandler) QuitJob(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cookie, err := c.Cookie("SESSION_ID")

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	var user quitJobDto

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ssoPayload := SsoPayload{
		Role:      "citizen",
		Initiator: user.Ucn,
		SessionId: cookie,
	}

	// change users role from employee to citizen in SSO
	if err := SendToService(c, ssoPayload, "http://sso_service:9090/user/employment/"+user.Ucn, true); err != nil {
		e.logger.Println(err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// payloadApr :=

	// if err := SendToService(c, payloadApr, "http://aprcroso:8005/adding-employee/request", false); err != nil {
	// 	e.logger.Println(err)
	// 	c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	// 	return
	// }

	var cv domain.CV

	cvColl := e.repo.GetCollection("employmentdb", "cvs")

	err = cvColl.FindOne(ctx, bson.M{"citizen_ucn": user.Ucn}).Decode(&cv)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

	if len(cv.WorkExperience) > 0 {
		lastIndex := len(cv.WorkExperience) - 1
		fieldPath := fmt.Sprintf("work_experience.%d", lastIndex)

		var workExperience string = time.Now().Month().String() +
			" " + strconv.Itoa(time.Now().Year()) + ")"

		_, err = cvColl.UpdateOne(ctx, bson.M{"citizen_ucn": user.Ucn},
			bson.M{"$set": bson.M{
				fieldPath: cv.WorkExperience[lastIndex] + workExperience,
			}})

		if err != nil {
			http.Error(c.Writer, err.Error(),
				http.StatusInternalServerError)
			e.logger.Println(err)
			return
		}
	}

	employeeColl := e.repo.GetCollection("employmentdb", "employees")

	_, err = employeeColl.DeleteOne(ctx, bson.D{{Key: "citizen_ucn", Value: user.Ucn}})

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		e.logger.Println(err)
		return
	}

}

type quitJobDto struct {
	Ucn string `bson:"ucn" json:"ucn"`
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
