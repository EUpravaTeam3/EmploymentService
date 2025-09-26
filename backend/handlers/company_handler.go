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
	/*ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	companyCollection := ch.repo.GetCollection(dbName, companyCollName)

	var company *domain.Company*/
	fmt.Println("COMMUNICATION EUPRAVA EMPLOYMENT")
	fmt.Println(c.Request.Body)

	/*if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := companyCollection.InsertOne(ctx, &company)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		ch.logger.Println(err)
		return
	}

	ch.logger.Printf("Documents ID: %v\n", result.InsertedID)
	e := json.NewEncoder(c.Writer)
	e.Encode(result)*/
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
	objID, err := primitive.ObjectIDFromHex(ownerUcn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid owner ucn"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var companyCollection = ch.repo.GetCollection(dbName, companyCollName)

	var company domain.Company
	err = companyCollection.FindOne(ctx, bson.M{"owner_ucn": objID}).Decode(&company)
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
