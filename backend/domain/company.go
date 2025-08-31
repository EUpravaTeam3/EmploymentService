package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CompanyName       string             `bson:"company_name" json:"company_name"`
	Industry          string             `bson:"industry" json:"industry"`
	Headquarters      string             `bson:"headquarters" json:"headquarters"`
	CompanyInfo       string             `bson:"company_info" json:"company_info"`
	NumberOfEmployees int                `bson:"number_of_employees" json:"number_of_employees"`
	IdNumber          int                `bson:"id_number" json:"id_number"`
	TaxIdNumber       int                `bson:"tax_id_number" json:"tax_id_number"`
	OwnerUcn          string             `bson:"owner_ucn" json:"owner_ucn"`
}

type Companies []*Company

func (c *Companies) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Company) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Company) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}
