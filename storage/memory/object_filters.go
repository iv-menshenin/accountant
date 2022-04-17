package memory

import (
	"fmt"
	"strings"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/storage"
)

func checkObjectFilter(object domain.Object, filter storage.FindObjectOption) bool {
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

func checkObjectAddress(object domain.Object, address string) bool {
	var addressA = strings.ToUpper(address)
	var addressB = strings.ToUpper(fmt.Sprintf("%s %s %s", object.City, object.Village, object.Street))
	return strings.Contains(addressB, addressA)
}

func checkObjectNumber(object domain.Object, number int) bool {
	return object.Number == number
}
