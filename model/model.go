package model

import "time"

type (
	UserIDType    string
	IDType        int64
	AttributeType string

	Attribute struct {
		Name string        `bson:"name"`
		Type AttributeType `bson:"type"`
	}
	ObjectType string
	Attributes struct {
		ObjectType ObjectType  `bson:"object_type"`
		Attributes []Attribute `bson:"attributes"`
	}

	Account struct {
		AccountID  IDType                 `bson:"account_id" json:"account_id"`
		Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
		Person     []Person               `bson:"persons" json:"persons"`
		Object     []Object               `bson:"objects" json:"objects"`
	}
	Person struct {
		PersonID   IDType                 `bson:"person_id" json:"person_id"`
		Name       string                 `bson:"name" json:"name"`
		Surname    string                 `bson:"surname" json:"surname"`
		PatName    string                 `bson:"pat_name" json:"pat_name"`
		DOB        time.Time              `bson:"dob,omitempty" json:"dob,omitempty"`
		Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
	}
	Object struct {
		ObjectID   IDType                 `bson:"object_id" json:"object_id"`
		Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
	}
)

const (
	AttributeTypeString  AttributeType = "String"
	AttributeTypeInteger AttributeType = "Integer"
	AttributeTypeDecimal AttributeType = "Decimal"
	AttributeTypeDate    AttributeType = "Date"

	ObjectTypeAccount ObjectType = "Account"
	ObjectTypePerson  ObjectType = "Person"
	ObjectTypeObject  ObjectType = "Object"
)
