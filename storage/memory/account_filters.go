package memory

import (
	"strings"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
)

func checkAccountFilter(account domain.Account, filter storage.FindAccountOption) bool {
	if filter.Account != nil {
		if !checkAccount(account, *filter.Account) {
			return false
		}
	}
	// todo fixme
	return true
}

func checkAccount(account domain.Account, pattern string) bool {
	return strings.Contains(account.Account, pattern)
}

func checkAccountFullName(account domain.Account, pattern string) bool {
	for _, person := range account.Persons {
		if checkPersonFullName(person, pattern) {
			return true
		}
	}
	return false
}

func checkAccountStreet(account domain.Account, pattern string) bool {
	var correct bool
	for _, object := range account.Objects {
		if correct = strings.Contains(object.Street, pattern); correct {
			break
		}
	}
	return correct
}

func checkAccountBuilding(account domain.Account, pattern int) bool {
	var correct bool
	for _, object := range account.Objects {
		if correct = object.Number == pattern; correct {
			break
		}
	}
	return correct
}

func checkAccountSumArea(account domain.Account, pattern float64) bool {
	var area float64
	for _, object := range account.Objects {
		area += object.Area
	}
	return area >= pattern
}
