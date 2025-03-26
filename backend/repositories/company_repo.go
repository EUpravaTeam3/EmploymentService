package repositories

import (
	"context"
	"employment-service/domain"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CompanyRepo struct {
	Cli    *mongo.Client
	logger *log.Logger
}

func NewCompanyRepo(ctx context.Context, logger *log.Logger) (*CompanyRepo, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://employment_db:27017/"))
	if err != nil {
		return nil, err
	}

	return &CompanyRepo{
		Cli:    client,
		logger: logger,
	}, nil
}

func (cr *CompanyRepo) GetCollection(collectionName string) *mongo.Collection {

	database := cr.Cli.Database("mongoDemo")
	collection := database.Collection(collectionName)
	return collection
}

func (cr *CompanyRepo) FindCompanyByJobId(jobId primitive.ObjectID) (domain.Company, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job domain.Job
	var company domain.Company

	jobCollection := cr.GetCollection("jobs")
	companyCollection := cr.GetCollection("companies")

	err := jobCollection.FindOne(ctx, bson.M{"_id": jobId}).Decode(&job)
	if err != nil {
		cr.logger.Println(err)
		return company, err
	}

	err = companyCollection.FindOne(ctx, bson.M{"_id": job.CompanyId}).Decode(&company)
	if err != nil {
		cr.logger.Println(err)
		return company, err
	}

	return company, nil
}
