package model

import "time"

type (
	AttributeType string
	Attribute     struct {
		Name string        `bson:"name"`
		Type AttributeType `bson:"type"`
	}
	ObjectType string
	Attributes struct {
		ObjectType ObjectType  `bson:"object_type"`
		Attributes []Attribute `bson:"attributes"`
	}

	Account struct {
		AccountID  int64                  `bson:"account_id"`
		Attributes map[string]interface{} `bson:"attributes"`
		Person     []Person               `bson:"persons"`
		Object     []Object               `bson:"objects"`
	}
	Person struct {
		Name       string                 `bson:"name"`
		Surname    string                 `bson:"surname"`
		PatName    string                 `bson:"pat_name"`
		DOB        time.Time              `bson:"dob,omitempty"`
		Attributes map[string]interface{} `bson:"attributes"`
	}
	Object struct {
		Attributes map[string]interface{} `bson:"attributes"`
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
