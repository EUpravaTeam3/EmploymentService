package handlers

import (
	"context"
	"employment-service/domain"
	"employment-service/repositories"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyHandler struct {
	logger *log.Logger
	repo   *repositories.CompanyRepo
}

func NewCompanyHandler(l *log.Logger, r *repositories.CompanyRepo) *CompanyHandler {
	return &CompanyHandler{l, r}
}

var companyDbName string = "employmentdb"
var companyCollName string = "companies"

func (ch *CompanyHandler) FindCompanyById(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var company, error = ch.repo.FindCompanyById(id)

	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	e := json.NewEncoder(c.Writer)
	e.Encode(company)
}

func (ch *CompanyHandler) CreateCompany(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	companyCollection := ch.repo.GetCollection(dbName, companyCollName)

	fmt.Println("COMMUNICATION EUPRAVA EMPLOYMENT")
	fmt.Println(c.Request.Body)

	var companyDTO RecievedCompany

	if err := c.ShouldBindJSON(&companyDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(companyDTO)

	var company domain.Company

	company.CompanyName = companyDTO.Name
	company.IdNumber = companyDTO.RegistrationNumber
	company.TaxIdNumber = companyDTO.PIB
	company.Status = companyDTO.CompanyStatus
	company.RegistrationDate = companyDTO.RegistrationDate
	company.OwnerUcn = companyDTO.OwnerUCN
	company.Address = companyDTO.Addresses
	company.WorkField = companyDTO.WorkFields

	_, errr := companyCollection.InsertOne(ctx, &company)
	if errr != nil {
		http.Error(c.Writer, errr.Error(),
			http.StatusInternalServerError)
		fmt.Println(errr)
		ch.logger.Println(errr)
		return
	}

	ssoPayload := SsoPayload{
		Role:      "employer",
		Initiator: company.OwnerUcn,
		SessionId: "eupravaDevelopment",
	}

	if err := SendToService(c, ssoPayload, "http://sso_service:9090/user/employment/"+company.OwnerUcn, true); err != nil {
		ch.logger.Println(err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"id":                 0,
		"name":               companyDTO.Name,
		"PIB":                companyDTO.PIB,
		"registrationNumber": companyDTO.RegistrationNumber,
		"registrationDate":   companyDTO.RegistrationDate.Format("2006-01-02"),
		"companyStatus":      companyDTO.CompanyStatus,
		"addresses":          companyDTO.Addresses,
		"employee":           companyDTO.Workers,
		"workFields":         []interface{}{},
		"ownerUcn":           companyDTO.OwnerUCN,
	}

	c.JSON(http.StatusCreated, response)
}

type RecievedCompany struct {
	ID                 interface{}     `bson:"_id,omitempty" json:"id,omitempty"`
	Name               string          `bson:"name" json:"name"`
	PIB                string          `bson:"PIB" json:"PIB"`
	RegistrationNumber string          `bson:"registrationNumber" json:"registrationNumber"`
	RegistrationDate   domain.DateOnly `bson:"registrationDate" json:"registrationDate"`
	CompanyStatus      string          `bson:"companyStatus" json:"companyStatus"`
	Addresses          []interface{}   `bson:"addresses" json:"addresses"`
	WorkFields         []interface{}   `bson:"workFields" json:"workFields"`
	Workers            []User          `bson:"workers" json:"worker"`
	OwnerUCN           string          `bson:"owner_ucn" json:"ownerUcn"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
}

func (ch *CompanyHandler) GetCompanies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var companies domain.Companies

	var companyCollection = ch.repo.GetCollection(dbName, companyCollName)

	companiesCursor, err := companyCollection.Find(ctx, bson.M{})
	if err != nil {
		ch.logger.Println(err)
		return
	}

	if err = companiesCursor.All(ctx, &companies); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		ch.logger.Fatal(err)
		return
	}

	err = companies.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		ch.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (ch *CompanyHandler) GetCompanyByOwnerUcn(c *gin.Context) {
	ownerUcn := c.Param("owner")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var companyCollection = ch.repo.GetCollection(dbName, companyCollName)

	var company domain.Company
	err := companyCollection.FindOne(ctx, bson.M{"owner_ucn": ownerUcn}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (ch *CompanyHandler) PostReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cookie, err := c.Cookie("SESSION_ID")

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusUnauthorized)
		ch.logger.Println(err)
		return
	}

	var review domain.ReviewOfCompany

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var employee domain.Employee

	emplColl := ch.repo.GetCollection("employmentdb", "employees")

	err = emplColl.FindOne(ctx, bson.M{"_id": review.EmployeeId}).Decode(&employee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := "http://sso_service:9090/user/check/" + employee.CitizenUCN

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.AddCookie(&http.Cookie{
		Name:  "SESSION_ID",
		Value: cookie,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	defer resp.Body.Close()

	var job domain.Job

	jobsColl := ch.repo.GetCollection("employmentdb", "jobs")

	err = jobsColl.FindOne(ctx, bson.M{"_id": employee.JobId}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	review.EmployerId = job.CompanyId

	reviewColl := ch.repo.GetCollection("employmentdb", "reviews")

	_, err = reviewColl.InsertOne(ctx, &review)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		ch.logger.Println(err)
		return
	}

	c.JSON(http.StatusOK, review)
}

func (ch *CompanyHandler) GetReviewsOfCompany(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := c.Param("company_id")

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		ch.logger.Fatal(err)
		return
	}

	var reviews domain.ReviewsOfCompany

	var reviewsCollection = ch.repo.GetCollection("employmentdb", "reviews")

	reviewsCursor, err := reviewsCollection.Find(ctx, bson.M{"employer_id": objectId})
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		ch.logger.Fatal(err)
		return
	}

	if err = reviewsCursor.All(ctx, &reviews); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		ch.logger.Fatal(err)
		return
	}

	err = reviews.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		ch.logger.Fatal("Unable to convert to json :", err)
		return
	}
}
