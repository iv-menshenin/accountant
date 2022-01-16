package model

import "github.com/iv-menshenin/accountant/model/uuid"

type (
	AuthData struct {
		JWT    string
		UserID uuid.UUID
		Roles  []string
	}

	Unauthorized struct{}
	Forbidden    struct{}
	NotFound     struct{}

	GetAccountQuery struct {
		ID uuid.UUID
	}
	PostAccountQuery struct {
		AccountData `json:",inline"`
	}
	PutAccountQuery struct {
		ID      uuid.UUID
		Account AccountData
	}
	DeleteAccountQuery struct {
		ID uuid.UUID
	}
	FindAccountsQuery struct {
		Account        *string
		Street         *string
		Building       *int
		PersonFullName *string
		SumArea        *float64
	}
)

func (u Unauthorized) Error() string {
	return "Authentication required"
}

func (f Forbidden) Error() string {
	return "You are not allowed this action"
}

func (n NotFound) Error() string {
	return "Object not found"
}
