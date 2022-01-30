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

func (q FindObjectOption) FillFromQuery(query model.FindObjectsQuery) {
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
