package domain

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diploma struct {
	Id              primitive.ObjectID `bson:"_id"`
	InstitutionId   uuid.UUID          `json:"institutionId"`
	InstitutionName string             `json:"institutionName"`
	InstitutionType string             `json:"institutionType"`
	AverageGrade    float64            `json:"averageGrade"`
	OwnerUCN        string             `json:"ucn"`
}
