package domain

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diploma struct {
	Id              primitive.ObjectID `bson:"_id"`
	InstitutionId   uuid.UUID          `bson:"institution_id"`
	InstitutionName string             `bson:"institution_name"`
	InstitutionType string             `bson:"institution_type"`
	AverageGrade    float64            `bson:"averageGrade"`
	OwnerUCN        string             `bson:"ucn"`
}
