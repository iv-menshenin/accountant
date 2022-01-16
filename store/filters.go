package store

import (
	"github.com/iv-menshenin/accountant/model"
	"strings"
)

type (
	FindAccountOption struct {
		Account        *string
		Street         *string
		Building       *int
		PersonFullName *string
		SumArea        *float64
	}
)

func (q FindAccountOption) FillFromQuery(query model.FindAccountsQuery) {
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

func checkAccountFilter(account model.Account, filter FindAccountOption) bool {
	if filter.Account != nil {
		if !checkAccount(account, *filter.Account) {
			return false
		}
	}
	if filter.PersonFullName != nil {
		if !checkAccountFullName(account, *filter.PersonFullName) {
			return false
		}
	}
	if filter.Street != nil {
		if !checkAccountStreet(account, *filter.Street) {
			return false
		}
	}
	if filter.Building != nil {
		if !checkAccountBuilding(account, *filter.Building) {
			return false
		}
	}
	if filter.SumArea != nil {
		if !checkAccountSumArea(account, *filter.SumArea) {
			return false
		}
	}
	return true
}

func checkAccount(account model.Account, pattern string) bool {
	return strings.Contains(account.Account, pattern)
}

func checkAccountFullName(account model.Account, pattern string) bool {
	var correct bool
	for _, person := range account.Person {
		correct = correct ||
			strings.Contains(person.Name, pattern) ||
			strings.Contains(person.Surname, pattern) ||
			strings.Contains(person.PatName, pattern)
		if correct {
			break
		}
	}
	return correct
}

func checkAccountStreet(account model.Account, pattern string) bool {
	var correct bool
	for _, object := range account.Object {
		if correct = strings.Contains(object.Street, pattern); correct {
			break
		}
	}
	return correct
}

func checkAccountBuilding(account model.Account, pattern int) bool {
	var correct bool
	for _, object := range account.Object {
		if correct = object.Number == pattern; correct {
			break
		}
	}
	return correct
}

func checkAccountSumArea(account model.Account, pattern float64) bool {
	var area float64
	for _, object := range account.Object {
		area += object.Area
	}
	return area >= pattern
}
