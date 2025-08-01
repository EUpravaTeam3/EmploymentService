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

type JobAdHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewJobAdHandler(l *log.Logger, r *repositories.EmploymentRepo) *JobAdHandler {
	return &JobAdHandler{l, r}
}

var jobAdCollName string = "jobAds"

func (j *JobAdHandler) GetJobAds(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var jobAds domain.JobAds

	var JobAdCollection = j.repo.GetCollection(dbName, jobAdCollName)

	JobAdCursor, err := JobAdCollection.Find(ctx, bson.M{})
	if err != nil {
		j.logger.Println(err)
		return
	}

	if err = JobAdCursor.All(ctx, &jobAds); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Fatal(err)
		return
	}

	err = jobAds.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (j *JobAdHandler) PostJobAd(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	jobAdCollection := j.repo.GetCollection(dbName, jobAdCollName)

	var jobAd *domain.JobAd

	if err := c.ShouldBindJSON(&jobAd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := jobAdCollection.InsertOne(ctx, &jobAd)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		j.logger.Println(err)
		return
	}
	j.logger.Printf("Documents ID: %v\n", result.InsertedID)
	e := json.NewEncoder(c.Writer)
	e.Encode(result)
}

func (j *JobAdHandler) GetJobAdById(c *gin.Context) {
	id := c.Param("jobad_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var jobAd domain.JobAd

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	jobAdCollection := j.repo.GetCollection(dbName, jobAdCollName)
	err = jobAdCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&jobAd)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Println(err)
	}
	jobAd.ToJSON(c.Writer)
}

func (j *JobAdHandler) EditJobAd(c *gin.Context) {
	id := c.Param("jobad_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var jobAd domain.JobAd

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BindJSON(&jobAd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobAdCollection := j.repo.GetCollection(dbName, jobAdCollName)

	jobAdCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{
		"$set": bson.M{
			"ad_title":        jobAd.AdTitle,
			"job_description": jobAd.JobDescription,
			"qualification":   jobAd.Qualification,
			"job_type":        jobAd.JobType,
		}})
}

func (j *JobAdHandler) DeleteJobAdById(c *gin.Context) {
	id := c.Param("jobad_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobAdCollection := j.repo.GetCollection(dbName, jobAdCollName)
	applicantsCollection := j.repo.GetCollection(dbName, ApplicantCollName)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, errApplicant := applicantsCollection.DeleteMany(ctx, bson.M{"job_ad_id": objectId})
	if errApplicant != nil {
		http.Error(c.Writer, errApplicant.Error(),
			http.StatusInternalServerError)
		j.logger.Println(errApplicant)
	}

	_, errJobAd := jobAdCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if errJobAd != nil {
		http.Error(c.Writer, errJobAd.Error(),
			http.StatusInternalServerError)
		j.logger.Println(errJobAd)
	}

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Println(err)
	}
}

func (e *JobAdHandler) JobadCORSMiddleware() gin.HandlerFunc {
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
