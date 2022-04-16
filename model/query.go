package model

import "github.com/iv-menshenin/accountant/model/uuid"

type (
	AuthData struct {
		JWT     string    `json:"jwt_token"`
		UserID  uuid.UUID `json:"user_id"`
		Context []string  `json:"context"`
		Refresh string    `json:"refresh_token"`
	}
	AuthQuery struct {
		Login    string `json:"login,omitempty"`
		Password string `json:"password,omitempty"`
	}
	RefreshTokenQuery struct {
		Token string `json:"token,omitempty"`
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
		Account *string
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

	PostObjectQuery struct {
		AccountID  uuid.UUID
		ObjectData ObjectData
	}
	GetObjectQuery struct {
		AccountID uuid.UUID
		ObjectID  uuid.UUID
	}
	PutObjectQuery struct {
		AccountID  uuid.UUID
		ObjectID   uuid.UUID
		ObjectData ObjectData
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
		Target TargetData
	}
	DeleteTargetQuery struct {
		TargetID uuid.UUID
	}
	FindTargetQuery struct {
		ShowClosed bool
		Period     *Period
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
