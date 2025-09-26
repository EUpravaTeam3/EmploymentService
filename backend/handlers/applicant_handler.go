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
)

type ApplicantHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewApplicantHandler(l *log.Logger, r *repositories.EmploymentRepo) *ApplicantHandler {
	return &ApplicantHandler{l, r}
}

var ApplicantCollName string = "applicants"

func (a *ApplicantHandler) GetApplicants(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var applicants domain.Applicants

	var applicantCollection = a.repo.GetCollection(dbName, ApplicantCollName)

	applicantsCursor, err := applicantCollection.Find(ctx, bson.M{})
	if err != nil {
		a.logger.Println(err)
		return
	}

	if err = applicantsCursor.All(ctx, &applicants); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		a.logger.Fatal(err)
		return
	}

	err = applicants.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *ApplicantHandler) GetApplicantsByJobad(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := c.Param("jobad_id")

	var applicants domain.Applicants
	var applicant domain.Applicant

	var applicantCollection = a.repo.GetCollection(dbName, ApplicantCollName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	applicantsCursor, err := applicantCollection.Find(ctx, bson.M{"job_ad_id": objectId})
	if err != nil {
		a.logger.Println(err)
		return
	}

	for applicantsCursor.Next(ctx) {
		if err := applicantsCursor.Decode(&applicant); err != nil {
			log.Fatalf("Failed to decode document: %v", err)
		}
		applicants = append(applicants, &applicant)
	}

	err = applicants.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *ApplicantHandler) PostApplicant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	applicantCollection := a.repo.GetCollection(dbName, ApplicantCollName)

	var applicant *domain.Applicant

	if err := c.ShouldBindJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cvColl := a.repo.GetCollection("employmentdb", "cvs")

	count, err := cvColl.CountDocuments(ctx, bson.M{"_id": applicant.CVId})

	if count == 0 {
		http.Error(c.Writer, "You don't have a CV!",
			http.StatusBadRequest)
		return
	}

	result, err := applicantCollection.InsertOne(ctx, &applicant)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		a.logger.Println(err)
		return
	}
	a.logger.Printf("Documents ID: %v\n", result.InsertedID)
	e := json.NewEncoder(c.Writer)
	e.Encode(result)
}

func (a *ApplicantHandler) GetApplicantById(c *gin.Context) {
	id := c.Param("applicant_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var applicant domain.Applicant

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	applicantCollection := a.repo.GetCollection(dbName, ApplicantCollName)
	err = applicantCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&applicant)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		a.logger.Println(err)
	}
	applicant.ToJSON(c.Writer)
}

func (a *ApplicantHandler) DeleteApplicantById(c *gin.Context) {
	id := c.Param("applicant_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	applicantCollection := a.repo.GetCollection(dbName, ApplicantCollName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	_, err = applicantCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objectId}})

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		a.logger.Println(err)
	}
}
