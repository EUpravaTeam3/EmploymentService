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

func (a *ApplicantHandler) GetApplicantsForCompany(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ucn := c.Param("ucn")

	pipeline := mongo.Pipeline{
		// 1. Match company by owner_ucn
		{{"$match", bson.D{{"owner_ucn", ucn}}}},

		// 2. Lookup jobs
		{{"$lookup", bson.D{
			{"from", "jobs"},
			{"localField", "_id"},
			{"foreignField", "company_id"},
			{"as", "jobs"},
		}}},
		{{"$unwind", "$jobs"}},

		// 3. Lookup jobAds
		{{"$lookup", bson.D{
			{"from", "jobAds"},
			{"localField", "jobs._id"},
			{"foreignField", "job_id"},
			{"as", "jobAds"},
		}}},
		{{"$unwind", "$jobAds"}},

		// 4. Lookup applicants
		{{"$lookup", bson.D{
			{"from", "applicants"},
			{"localField", "jobAds._id"},
			{"foreignField", "job_ad_id"},
			{"as", "applicants"},
		}}},
		{{"$unwind", "$applicants"}},

		// 5. Lookup CVs
		{{"$lookup", bson.D{
			{"from", "cvs"},
			{"localField", "applicants.cv_id"},
			{"foreignField", "_id"},
			{"as", "cv"},
		}}},
		{{"$unwind", "$cv"}},

		// 6. Project only needed fields
		{{"$project", bson.D{
			{"position_name", "$jobs.position_name"},
			{"ad_title", "$jobAds.ad_title"},
			{"job_ad_id", "$jobAds._id"},
			{"citizen_ucn", "$cv.citizen_ucn"},
			{"company_owner_ucn", ucn},
			{"name", "$cv.name"},
			{"email", "$cv.email"},
			{"description", "$cv.description"},
			{"work_experience", "$cv.work_experience"},
			{"education", "$cv.education"},
		}}},
	}

	companyColl := a.repo.GetCollection("employmentdb", "companies")

	cursor, err := companyColl.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "aggregation failed", "details": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read cursor", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *ApplicantHandler) PostApplicant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ucn := c.Param("ucn")

	applicantCollection := a.repo.GetCollection(dbName, ApplicantCollName)

	var applicant *domain.Applicant

	if err := c.ShouldBindJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("AAPLICANT")
	fmt.Println(applicant.JobAdId)
	fmt.Println(ucn)
	fmt.Println(applicant.CVId)

	cvColl := a.repo.GetCollection("employmentdb", "cvs")

	count, err := cvColl.CountDocuments(ctx, bson.M{"citizen_ucn": ucn})

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(c.Writer, "You don't have a CV!",
			http.StatusBadRequest)
		return
	}

	if count == 0 {
		http.Error(c.Writer, "You don't have a CV!",
			http.StatusBadRequest)
		return
	}

	var cv domain.CV

	err = cvColl.FindOne(ctx, bson.M{"citizen_ucn": ucn}).Decode(&cv)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		return
	}

	applicant.CVId = cv.Id

	count, err = applicantCollection.CountDocuments(ctx, bson.M{"cv_id": applicant.CVId,
		"job_ad_id": applicant.JobAdId})

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		a.logger.Println(err)
		return
	}

	if count > 0 {
		http.Error(c.Writer, "You already applied for this position!",
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

func (a *ApplicantHandler) GetApplicantByUcn(c *gin.Context) {
	ucn := c.Param("ucn")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	applicantsCollection := a.repo.GetCollection("employmentdb", "applicants")

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "cvs"},
			{"localField", "cv_id"},
			{"foreignField", "_id"},
			{"as", "cv"},
		}}},
		{{"$unwind", "$cv"}},
		{{"$match", bson.D{
			{"cv.citizen_ucn", ucn},
		}}},
		{{"$lookup", bson.D{
			{"from", "jobAds"},
			{"localField", "job_ad_id"},
			{"foreignField", "_id"},
			{"as", "jobAd"},
		}}},
		{{"$unwind", "$jobAd"}},
		{{"$project", bson.D{
			{"_id", 0},
			{"applicant_id", "$_id"},
			{"job_ad_id", "$job_ad_id"},
			{"cv_id", "$cv_id"},
			{"job_ad_title", "$jobAd.ad_title"},
		}}},
	}

	cursor, err := applicantsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal("Aggregate error:", err)
	}
	defer cursor.Close(ctx)

	var results []ApplicantWithJob
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatal("Decode error:", err)
	}

	c.JSON(http.StatusOK, results)

}

func (a *ApplicantHandler) DeleteApplicantById(c *gin.Context) {
	id := c.Param("applicant_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println(id)
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

type ApplicantWithJob struct {
	ApplicantID string `bson:"applicant_id" json:"applicant_id"`
	JobAdID     string `bson:"job_ad_id" json:"job_ad_id"`
	CVID        string `bson:"cv_id" json:"cv_id"`
	JobAdTitle  string `bson:"job_ad_title" json:"job_ad_title"`
}
