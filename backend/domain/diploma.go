package domain

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diploma struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	InstitutionId   uuid.UUID          `bson:"institution_id"`
	InstitutionName string             `bson:"institution_name"`
	InstitutionType string             `bson:"institution_type"`
	AverageGrade    float64            `bson:"average_grade"`
	OwnerUCN        string             `bson:"ucn"`
}
