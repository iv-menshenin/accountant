package store

import (
	"fmt"
	"strings"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	FindObjectOption struct {
		AccountID *uuid.UUID
		Address   *string
		Number    *int
	}
)

func checkObjectFilter(object model.Object, filter FindObjectOption) bool {
	if filter.Address != nil {
		if !checkObjectAddress(object, *filter.Address) {
			return false
		}
	}
	if filter.Number != nil {
		if !checkObjectNumber(object, *filter.Number) {
			return false
		}
	}
	return true
}

func checkObjectAddress(object model.Object, address string) bool {
	var addressA = strings.ToUpper(address)
	var addressB = strings.ToUpper(fmt.Sprintf("%s %s %s", object.City, object.Village, object.Street))
	return strings.Contains(addressB, addressA)
}

func checkObjectNumber(object model.Object, number int) bool {
	return object.Number == number
}
