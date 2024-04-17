package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employer struct {
	Id                primitive.ObjectID `bson:"_id"`
	CompanyName       string             `bson:"company_name"`
	Industry          string             `bson:"industry"`
	Headquarters      string             `bson:"headquarters"`
	CompanyInfo       string             `bson:"company_info"`
	NumberOfEmployees int                `bson:"number_of_employees"`
	IdNumber          int                `bson:"id_number"`
	TaxIdNumber       int                `bson:"tax_id_number"`
	OwnerUcn          string             `bson:"owner_ucn"`
}
