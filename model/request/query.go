package request

import (
	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	GetAccountQuery struct {
		ID uuid.UUID
	}
	PostAccountQuery struct {
		AccountData domain.AccountData
	}
	PutAccountQuery struct {
		ID      uuid.UUID
		Account domain.AccountData
	}
	DeleteAccountQuery struct {
		ID uuid.UUID
	}
	FindAccountsQuery struct {
		Account *string
	}

	PostPersonQuery struct {
		AccountID  uuid.UUID
		PersonData domain.PersonData
	}
	GetPersonQuery struct {
		AccountID uuid.UUID
		PersonID  uuid.UUID
	}
	PutPersonQuery struct {
		AccountID  uuid.UUID
		PersonID   uuid.UUID
		PersonData domain.PersonData
	}
	DeletePersonQuery struct {
		AccountID uuid.UUID
		PersonID  uuid.UUID
	}
	FindPersonsQuery struct {
		AccountID      *uuid.UUID
		PersonFullName *string
	}

	PostObjectQuery struct {
		AccountID  uuid.UUID
		ObjectData domain.ObjectData
	}
	GetObjectQuery struct {
		AccountID uuid.UUID
		ObjectID  uuid.UUID
	}
	PutObjectQuery struct {
		AccountID  uuid.UUID
		ObjectID   uuid.UUID
		ObjectData domain.ObjectData
	}
	DeleteObjectQuery struct {
		AccountID uuid.UUID
		ObjectID  uuid.UUID
	}
	FindObjectsQuery struct {
		AccountID *uuid.UUID
		Address   *string
		Number    *int
	}

	GetTargetQuery struct {
		TargetID uuid.UUID
	}
	PostTargetQuery struct {
		Type   string
		Target domain.TargetData
	}
	DeleteTargetQuery struct {
		TargetID uuid.UUID
	}
	FindTargetQuery struct {
		ShowClosed bool
		Period     *domain.Period
	}
)