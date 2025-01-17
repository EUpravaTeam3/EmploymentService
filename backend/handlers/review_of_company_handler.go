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

type ReviewOfCompanyHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewReviewOfCompanyHandler(l *log.Logger, r *repositories.EmploymentRepo) *ReviewOfCompanyHandler {
	return &ReviewOfCompanyHandler{l, r}
}

var ReviewOfCompanyCollName string = "reviewsOfCompany"

func (rw *ReviewOfCompanyHandler) GetReviewsOfCompany(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("company_id")

	var review domain.ReviewOfCompany
	var reviews domain.ReviewsOfCompany

	var reviewCollection = rw.repo.GetCollection(dbName, ReviewOfCompanyCollName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	reviewsCursor, err := reviewCollection.Find(ctx, bson.M{"employer_id": objectId})
	if err != nil {
		rw.logger.Println(err)
		return
	}

	for reviewsCursor.Next(ctx) {
		if err := reviewsCursor.Decode(&review); err != nil {
			log.Fatalf("Failed to decode document: %v", err)
		}
		reviews = append(reviews, &review)
	}

	err = reviews.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		rw.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (rw *ReviewOfCompanyHandler) PostReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reviewCollection := rw.repo.GetCollection(dbName, ReviewOfCompanyCollName)

	var review *domain.ReviewOfCompany

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := reviewCollection.InsertOne(ctx, &review)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		rw.logger.Println(err)
		return
	}
	rw.logger.Printf("Documents ID: %v\n", result.InsertedID)
	e := json.NewEncoder(c.Writer)
	e.Encode(result)
}

func (rw *ReviewOfCompanyHandler) GetReviewById(c *gin.Context) {
	id := c.Param("review_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var review domain.ReviewOfCompany

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	reviewCollection := rw.repo.GetCollection(dbName, ReviewOfCompanyCollName)
	err = reviewCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&review)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		rw.logger.Println(err)
	}
	review.ToJSON(c.Writer)
}

func (rw *ReviewOfCompanyHandler) EditReview(c *gin.Context) {
	id := c.Param("review_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var review domain.ReviewOfCompany

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviewCollection := rw.repo.GetCollection(dbName, ReviewOfCompanyCollName)

	reviewCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{
		"$set": bson.M{
			"description": review.Description,
			"rating":      review.Rating,
		}})
}

func (rw *ReviewOfCompanyHandler) DeleteReviewById(c *gin.Context) {
	id := c.Param("review_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reviewCollection := rw.repo.GetCollection(dbName, ReviewOfCompanyCollName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	_, err = reviewCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objectId}})

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		rw.logger.Println(err)
	}
}
