package handlers

import (
	"employment-service/domain"
	"employment-service/repositories"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type EmploymentHandler struct {
	logger *log.Logger
	repo   *repositories.EmploymentRepo
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func NewEmploymentHandler(l *log.Logger, r *repositories.EmploymentRepo) *EmploymentHandler {
	return &EmploymentHandler{l, r}
}

func (e *EmploymentHandler) GenerateDiploma(c *gin.Context) {

	var diploma *domain.Diploma

	if err := c.BindJSON(&diploma); err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := e.repo.AddDiplomaToCV(diploma); err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diploma)

}

func (p *EmploymentHandler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
