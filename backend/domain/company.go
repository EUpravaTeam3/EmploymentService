package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CompanyName      string             `bson:"company_name" json:"company_name"`
	Status           string             `bson:"status" json:"status"`
	IdNumber         string             `bson:"id_number" json:"id_number"`
	TaxIdNumber      string             `bson:"tax_id_number" json:"tax_id_number"`
	OwnerUcn         string             `bson:"owner_ucn" json:"owner_ucn"`
	RegistrationDate DateOnly           `bson:"registration_date" json:"registration_date"`
	Address          []interface{}      `bson:"address" json:"address"`
	WorkField        []interface{}      `bson:"work_field" json:"work_field"`
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
