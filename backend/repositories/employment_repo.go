package repositories

import (
	"context"
	"employment-service/domain"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type EmploymentRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*EmploymentRepo, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://employment_db:27017/"))
	if err != nil {
		return nil, err
	}

	return &EmploymentRepo{
		cli:    client,
		logger: logger,
	}, nil
}

func (er *EmploymentRepo) AddDiplomaToCV(diploma *domain.Diploma) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	diplomaCollection := er.getCollection("cvs")

	filter := bson.M{"citizen_ucn": diploma.citizen_ucn}

	var existingCv bson.M

	if err := diplomaCollection.FindOne(ctx, filter).Decode(&existingCv); err != nil {

		if err == mongo.ErrNoDocuments {

			var newCv domain.CV
			newCv.CitizenUCN = diploma.citizen_ucn
			newCv.Description = ""
			newCV.WorkExperience = []string
			newCV.Education = []Diploma
			newCV.Education[0] = diploma

			result, err := collection.InsertOne(ctx, &newCV)
			if err != nil {
				fmt.Println(err)
				ar.logger.Println(err)
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

	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Print(err)
		return err
	}

	er.logger.Printf("result: %v\n", updateResult)

}

func getCollection(collectionName string) *mongo.Collection {

	database := er.cli.Database("mongoDemo")
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
func (pr *EmploymentRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *EmploymentRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if there's no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}
