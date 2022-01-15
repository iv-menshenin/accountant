package store

import (
	"github.com/iv-menshenin/accountant/model"
	"strings"
)

type (
	FindAccountOption struct {
		Street         *string
		Building       *int
		PersonFullName *string
	}
)

func checkAccountFilter(account model.Account, filter FindAccountOption) bool {
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
	return true
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
