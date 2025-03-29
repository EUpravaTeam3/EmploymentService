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

type EmploymentRepo struct {
	Cli    *mongo.Client
	logger *log.Logger
}

func NewEmp(ctx context.Context, logger *log.Logger) (*EmploymentRepo, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://employment_db:27017/"))
	if err != nil {
		return nil, err
	}

	return &EmploymentRepo{
		Cli:    client,
		logger: logger,
	}, nil
}

func (er *EmploymentRepo) GetCollection(dbName, collectionName string) *mongo.Collection {
	return er.Cli.Database(dbName).Collection(collectionName)
}

func (er *EmploymentRepo) AddDiplomaToCV(diploma *domain.Diploma) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cvCollection := er.getCollection("cvs")

	filter := bson.M{"citizen_ucn": diploma.OwnerUCN}

	var existingCv bson.M

	if err := cvCollection.FindOne(ctx, filter).Decode(&existingCv); err != nil {

		if err == mongo.ErrNoDocuments {

			var newCv domain.CV
			newCv.CitizenUCN = diploma.OwnerUCN
			newCv.Description = ""
			newCv.WorkExperience = []string{}
			newCv.Education = []domain.Diploma{}
			newCv.Education = append(newCv.Education, *diploma)

			result, err := cvCollection.InsertOne(ctx, &newCv)
			if err != nil {
				fmt.Println(err)
				er.logger.Println(err)
				return err
			}

			er.logger.Printf("Documents ID: %v\n", result.InsertedID)
			return nil

		} else {
			log.Fatal(err)
			return err
		}
	}

	update := bson.M{"$push": bson.M{"education": diploma}}

	updateResult, err := cvCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Print(err)
		return err
	}

	er.logger.Printf("result: %v\n", updateResult)

	return nil

}

func (er *EmploymentRepo) AddEmployeeToCompany(employee domain.Employee) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	employeeCollection := er.getCollection("employees")

	_, err := employeeCollection.InsertOne(ctx, &employee)

	if err != nil {
		return err
	}

	return nil

}

func (er *EmploymentRepo) FindCompanyByJobId(jobId primitive.ObjectID) (domain.Company, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job domain.Job
	var company domain.Company

	jobCollection := er.getCollection("jobs")
	companyCollection := er.getCollection("companies")

	err := jobCollection.FindOne(ctx, bson.M{"_id": jobId}).Decode(&job)
	if err != nil {
		er.logger.Println(err)
		return company, err
	}

	err = companyCollection.FindOne(ctx, bson.M{"_id": job.CompanyId}).Decode(&company)
	if err != nil {
		er.logger.Println(err)
		return company, err
	}

	return company, nil
}

func (er *EmploymentRepo) InsertJobAd(jobAd *domain.JobAd) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobAdCollection := er.getCollection("jobAds")

	_, err := jobAdCollection.InsertOne(ctx, &jobAd)

	if err != nil {
		return err
	}

	return nil
}

func (er *EmploymentRepo) getCollection(collectionName string) *mongo.Collection {

	database := er.Cli.Database("mongoDemo")
	collection := database.Collection(collectionName)
	return collection
}

func DBinstance() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://employment_db:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()

// Disconnect from database
func (pr *EmploymentRepo) DisconnectEmp(ctx context.Context) error {
	err := pr.Cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *EmploymentRepo) PingEmp() {
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
