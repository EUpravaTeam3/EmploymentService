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
)

type UserHandler struct {
	logger *log.Logger
	repo   *repositories.UserRepo
}

func NewUserHandler(l *log.Logger, r *repositories.UserRepo) *UserHandler {
	return &UserHandler{l, r}
}

var userDbName string = "employmentdb"
var userCollName string = "users"

func (ch *UserHandler) FindUserByUcn(ucn string) (domain.User, error) {

	var user, error = ch.repo.FindUserByUcn(ucn)

	if error != nil {
		return user, error
	}

	return user, nil
}

func (ch *UserHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userCollection := ch.repo.GetCollection(userCollName)

	var user *domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := userCollection.InsertOne(ctx, &user)
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println(err)
		ch.logger.Println(err)
		return
	}

	ch.logger.Printf("Documents ID: %v\n", result.InsertedID)
	e := json.NewEncoder(c.Writer)
	e.Encode(result)
}
