package handlers

import (
	"context"
	"employment-service/domain"
	"employment-service/repositories"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CvHandler struct {
	logger *log.Logger
	repo   *repositories.CompanyRepo
}

func NewCvHandler(l *log.Logger, r *repositories.CompanyRepo) *CvHandler {
	return &CvHandler{l, r}
}

var cvDbName string = "employmentdb"
var cvCollName string = "cvs"

func (cvh *CvHandler) FindCvByUcn(c *gin.Context) {
	ucn := c.Param("ucn")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cv domain.CV

	cvCollection := cvh.repo.GetCollection(cvDbName, cvCollName)
	err := cvCollection.FindOne(ctx, bson.M{"citizen_ucn": ucn}).Decode(&cv)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		cvh.logger.Println(err)
	}
	cv.ToJSON(c.Writer)
}

func (cvh *CvHandler) PostCv(c *gin.Context) {
	ucn := c.Param("ucn")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cv domain.CV

	if err := c.BindJSON(&cv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cvCollection := cvh.repo.GetCollection(cvDbName, cvCollName)
	count, err := cvCollection.CountDocuments(ctx, bson.M{"citizen_ucn": ucn})
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		cvh.logger.Println(err)
	}

	if count > 0 {

		filter := bson.M{"citizen_ucn": ucn}
		update := bson.M{
			"$set": bson.M{
				"description":     cv.Description,
				"work_experience": cv.WorkExperience,
			},
		}
		_, err := cvCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(c.Writer, err.Error(),
				http.StatusInternalServerError)
			cvh.logger.Println(err)
		}

	} else {
		_, err := cvCollection.InsertOne(ctx, &cv)
		if err != nil {
			http.Error(c.Writer, err.Error(),
				http.StatusInternalServerError)
			fmt.Println(err)
			cvh.logger.Println(err)
			return
		}
	}

	cv.ToJSON(c.Writer)
}
