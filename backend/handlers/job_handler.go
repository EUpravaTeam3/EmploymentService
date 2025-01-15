package handlers

import (
	"context"
	"employment-service/domain"
	"employment-service/repositories"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewJobHandler(l *log.Logger, r *repositories.EmploymentRepo) *JobHandler {
	return &JobHandler{l, r}
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var dbName string = "employmentdb"
var collName string = "jobs"

func (j *JobHandler) GetJobs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var jobs domain.Jobs

	var jobCollection = j.repo.GetCollection(dbName, collName)

	jobsCursor, err := jobCollection.Find(ctx, bson.M{})
	if err != nil {
		j.logger.Println(err)
		return
	}

	if err = jobsCursor.All(ctx, &jobs); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Fatal(err)
		return
	}

	err = jobs.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (j *JobHandler) PostJob(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	jobCollection := j.repo.GetCollection(dbName, collName)

	var job *domain.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := jobCollection.InsertOne(ctx, &job)
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

func (j *JobHandler) GetJobById(c *gin.Context) {
	id := c.Param("job_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job domain.Job

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	jobCollection := j.repo.GetCollection(dbName, collName)
	err = jobCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&job)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Println(err)
	}
	job.ToJSON(c.Writer)
}

func (j *JobHandler) EditJob(c *gin.Context) {
	id := c.Param("job_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job domain.Job

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobCollection := j.repo.GetCollection(dbName, collName)

	jobCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{
		"$set": bson.M{
			"position_name":     job.PoistionName,
			"pay":               job.Pay,
			"num_of_employees":  job.NumOfEmployees,
			"Employee_capacity": job.EmployeeCapacity,
		}})
}

func (j *JobHandler) DeleteJobById(c *gin.Context) {
	id := c.Param("job_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobCollection := j.repo.GetCollection(dbName, collName)

	_, err := jobCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		j.logger.Println(err)
	}
}
