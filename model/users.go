package model

import "github.com/iv-menshenin/accountant/model/uuid"

type (
	User struct {
		UUID     uuid.UUID `bson:"user_id" json:"user_id"`
		UserName string    `bson:"user_name" json:"user_name"`
		Context  []string  `bson:"context" json:"context"`
	}
)
