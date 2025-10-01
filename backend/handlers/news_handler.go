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

type NewsHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

func NewNewsHandler(l *log.Logger, r *repositories.EmploymentRepo) *NewsHandler {
	return &NewsHandler{l, r}
}

var newsCollName string = "news"

func (n *NewsHandler) GetNews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var news domain.AllNews

	var newsCollection = n.repo.GetCollection(dbName, newsCollName)

	newsCursor, err := newsCollection.Find(ctx, bson.M{})
	if err != nil {
		n.logger.Println(err)
		return
	}

	if err = newsCursor.All(ctx, &news); err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		n.logger.Fatal(err)
		return
	}

	err = news.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		n.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (n *NewsHandler) PostNews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	newsCollection := n.repo.GetCollection(dbName, newsCollName)

	var news *domain.News

	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := newsCollection.InsertOne(ctx, &news)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		n.logger.Println(err)
		return
	}
	n.logger.Printf("Documents ID: %v\n", result.InsertedID)
	e := json.NewEncoder(c.Writer)
	e.Encode(result)
}

func (n *NewsHandler) GetNewsById(c *gin.Context) {
	id := c.Param("news_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var news domain.News

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	newsCollection := n.repo.GetCollection(dbName, newsCollName)
	err = newsCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&news)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		n.logger.Println(err)
	}
	news.ToJSON(c.Writer)
}

func (n *NewsHandler) EditNews(c *gin.Context) {
	id := c.Param("news_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var news domain.News

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newsCollection := n.repo.GetCollection(dbName, newsCollName)

	newsCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{
		"$set": bson.M{
			"title": news.Title,
		}})
}

func (n *NewsHandler) DeleteNewsById(c *gin.Context) {
	id := c.Param("news_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, _ := primitive.ObjectIDFromHex(id)

	newsCollection := n.repo.GetCollection(dbName, newsCollName)

	_, err := newsCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objectId}})

	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		n.logger.Println(err)
	}
}
