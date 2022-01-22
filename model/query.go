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
		AccountData AccountData
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

	PostPersonQuery struct {
		AccountID  uuid.UUID
		PersonData PersonData
	}
	GetPersonQuery struct {
		AccountID uuid.UUID
		PersonID  uuid.UUID
	}
	PutPersonQuery struct {
		AccountID  uuid.UUID
		PersonID   uuid.UUID
		PersonData PersonData
	}
	DeletePersonQuery struct {
		AccountID uuid.UUID
		PersonID  uuid.UUID
	}
	FindPersonsQuery struct {
		AccountID      *uuid.UUID
		PersonFullName *string
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
