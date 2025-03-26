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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (ch *CompanyHandler) FindCompanyByJobId(id primitive.ObjectID) (domain.Company, error) {

	var company, error = ch.repo.FindCompanyByJobId(id)

	if error != nil {
		return company, error
	}

	return company, nil
}

func (ch *CompanyHandler) CreateCompany(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	companyCollection := ch.repo.GetCollection(companyCollName)

	var company *domain.Company

	if err := c.ShouldBindJSON(&company); err != nil {
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
	e.Encode(result)
}
