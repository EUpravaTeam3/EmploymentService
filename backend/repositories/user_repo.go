package repositories

import (
	"context"
	"employment-service/domain"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	Cli    *mongo.Client
	logger *log.Logger
}

func NewUserRepo(ctx context.Context, logger *log.Logger) (*UserRepo, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://employment_db:27017/"))
	if err != nil {
		return nil, err
	}

	return &UserRepo{
		Cli:    client,
		logger: logger,
	}, nil
}

func (ur *UserRepo) GetCollection(collectionName string) *mongo.Collection {

	database := ur.Cli.Database("mongoDemo")
	collection := database.Collection(collectionName)
	return collection
}

func (ur *UserRepo) FindUserByUcn(ucn string) (domain.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User

	jobCollection := ur.GetCollection("users")

	err := jobCollection.FindOne(ctx, bson.M{"ucn": ucn}).Decode(&user)
	if err != nil {
		ur.logger.Println(err)
		return user, err
	}

	return user, nil
}
