package domain

import (
	"fmt"

	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	UserInfo struct {
		ID          uuid.UUID    `json:"user_id" bson:"user_id"`
		Name        string       `json:"name" bson:"name"`
		Surname     string       `json:"surname" bson:"surname"`
		EMail       string       `json:"e_mail" bson:"e_mail"`
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

func StringsToPermissions(perm []string) ([]Permission, error) {
	var result = make([]Permission, 0, len(perm))
	for _, s := range perm {
		p := Permission(s)
		if !checkAllowedPermission(p) {
			return nil, fmt.Errorf("permission %s is not allowed", s)
		}
		result = append(result, p)
	}
	return result, nil
}

func checkAllowedPermission(perm Permission) bool {
	switch perm {

	case PermView, PermWrite, PermFinance, PermAdmin:
		return true

	default:
		return false
	}
}
