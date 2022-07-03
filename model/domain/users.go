package domain

import "github.com/iv-menshenin/accountant/utils/uuid"

type (
	UserInfo struct {
		ID          uuid.UUID    `json:"user_id" bson:"user_id"`
		Name        string       `json:"name" bson:"name"`
		Surname     string       `json:"surname" bson:"surname"`
		Permissions []Permission `json:"permissions" bson:"permissions"`
	}
	UserIdentity struct {
		Login    string `json:"login" bson:"login"`
		Password string `json:"password" bson:"password"`
	}
	Permission string
)

const (
	PermView    Permission = "view"
	PermWrite   Permission = "write"
	PermFinance Permission = "finance"
	PermAdmin   Permission = "admin"
)
