package repositories

import (
	"context"
	"employment-service/domain"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func (cr *CompanyRepo) GetCollection(dbName string, collectionName string) *mongo.Collection {

	return cr.Cli.Database(dbName).Collection(collectionName)
}

func (cr *CompanyRepo) FindCompanyById(companyId primitive.ObjectID) (domain.Company, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var company domain.Company

	companyCollection := cr.GetCollection("employmentdb", "companies")

	var err = companyCollection.FindOne(ctx, bson.M{"_id": companyId}).Decode(&company)
	if err != nil {
		cr.logger.Println(err)
		return company, err
	}

	return company, nil
}

// Disconnect from database
func (pr *CompanyRepo) DisconnectComp(ctx context.Context) error {
	err := pr.Cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *CompanyRepo) PingComp() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if there's no error, connection is established
	err := pr.Cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.Cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}
