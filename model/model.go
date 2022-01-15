package model

import (
	"time"

	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	AccountData struct {
		Comment       string     `bson:"comment" json:"comment"`
		AgreementNum  string     `bson:"agreement" json:"agreement"`
		AgreementDate *time.Time `bson:"agreement_date,omitempty" json:"agreement_date,omitempty"`
	}
	Account struct {
		AccountID   uuid.UUID `bson:"account_id" json:"account_id"`
		Person      []Person  `bson:"persons" json:"persons"`
		Object      []Object  `bson:"objects" json:"objects"`
		AccountData `bson:",inline" json:",inline"`
	}
	PersonData struct {
		Name    string    `bson:"name" json:"name"`
		Surname string    `bson:"surname" json:"surname"`
		PatName string    `bson:"pat_name" json:"pat_name"`
		DOB     time.Time `bson:"dob,omitempty" json:"dob,omitempty"`
	}
	Person struct {
		PersonID   uuid.UUID `bson:"person_id" json:"person_id"`
		PersonData `bson:",inline" json:",inline"`
	}
	ObjectData struct {
		PostalCode string `bson:"postal_code" json:"postal_code"`
		City       string `bson:"city" json:"city"`
		Village    string `bson:"village" json:"village"`
		Street     string `bson:"street" json:"street"`
		Number     int    `bson:"number" json:"number"`
	}
	Object struct {
		ObjectID   uuid.UUID `bson:"object_id" json:"object_id"`
		ObjectData `bson:",inline" json:",inline"`
	}
)
