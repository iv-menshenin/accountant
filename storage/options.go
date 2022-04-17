package storage

import (
	"github.com/iv-menshenin/accountant/model/request"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

// todo move to storage

type (
	FindAccountOption struct {
		Account *string
		// hidden options
		Address *string
		Number  *int
	}
	FindObjectOption struct {
		AccountID *uuid.UUID
		Address   *string
		Number    *int
	}
	FindPersonOption struct {
		AccountID *uuid.UUID
	}
	FindTargetOption struct {
		ShowClosed bool
		Year       int
		Month      int
	}
)

func (q *FindAccountOption) FillFromQuery(query request.FindAccountsQuery) {
	if query.Account != nil {
		q.Account = query.Account
	}
}

func (q *FindObjectOption) FillFromQuery(query request.FindObjectsQuery) {
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

func (q *FindPersonOption) FillFromQuery(query request.FindPersonsQuery) {
	if query.AccountID != nil {
		q.AccountID = query.AccountID
	}
}
