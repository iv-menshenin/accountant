package memory

import (
	"strings"

	"github.com/iv-menshenin/accountant/model"
)

func checkAccountFilter(account model.Account, filter model.FindAccountOption) bool {
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
	for _, person := range account.Persons {
		if checkPersonFullName(person, pattern) {
			return true
		}
	}
	return false
}

func checkAccountStreet(account model.Account, pattern string) bool {
	var correct bool
	for _, object := range account.Objects {
		if correct = strings.Contains(object.Street, pattern); correct {
			break
		}
	}
	return correct
}

func checkAccountBuilding(account model.Account, pattern int) bool {
	var correct bool
	for _, object := range account.Objects {
		if correct = object.Number == pattern; correct {
			break
		}
	}
	return correct
}

func checkAccountSumArea(account model.Account, pattern float64) bool {
	var area float64
	for _, object := range account.Objects {
		area += object.Area
	}
	return area >= pattern
}
