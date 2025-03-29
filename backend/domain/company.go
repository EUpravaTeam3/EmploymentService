package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	CompanyName       string             `bson:"company_name"`
	Industry          string             `bson:"industry"`
	Headquarters      string             `bson:"headquarters"`
	CompanyInfo       string             `bson:"company_info"`
	NumberOfEmployees int                `bson:"number_of_employees"`
	IdNumber          int                `bson:"id_number"`
	TaxIdNumber       int                `bson:"tax_id_number"`
	OwnerUcn          string             `bson:"owner_ucn"`
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
