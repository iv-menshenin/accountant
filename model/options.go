package model

import (
	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	FindAccountOption struct {
		Account        *string
		Street         *string
		Building       *int
		PersonFullName *string
		SumArea        *float64
	}
	FindObjectOption struct {
		AccountID *uuid.UUID
		Address   *string
		Number    *int
	}
	FindPersonOption struct {
		AccountID      *uuid.UUID
		PersonFullName *string
	}
)

func (q *FindAccountOption) FillFromQuery(query FindAccountsQuery) {
	if query.Building != nil {
		q.Building = query.Building
	}
	if query.Street != nil {
		q.Street = query.Street
	}
	if query.SumArea != nil {
		q.SumArea = query.SumArea
	}
	if query.PersonFullName != nil {
		q.PersonFullName = query.PersonFullName
	}
	if query.Account != nil {
		q.Account = query.Account
	}
}

func (q *FindObjectOption) FillFromQuery(query FindObjectsQuery) {
	if query.Address != nil {
		q.Address = query.Address
	}
	if query.Number != nil {
		q.Number = query.Number
	}
	if query.AccountID != nil {
		q.AccountID = query.AccountID
	}
}

func (q *FindPersonOption) FillFromQuery(query FindPersonsQuery) {
	if query.PersonFullName != nil {
		q.PersonFullName = query.PersonFullName
	}
	if query.AccountID != nil {
		q.AccountID = query.AccountID
	}
}
