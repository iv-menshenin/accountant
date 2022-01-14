package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_SaveAccountQuery(t *testing.T) {
	var acc = Account{
		AccountID: IDType{},
		AccountData: AccountData{
			Attributes: map[string]interface{}{
				"foo": "bar",
			},
			Person: []Person{},
			Object: []Object{},
		},
	}
	var aq = SaveAccountQuery{
		AccountData: acc.AccountData,
	}
	data, err := json.Marshal(aq)
	if err != nil {
		t.Error(err)
	}
	var acc2 Account
	if err = json.Unmarshal(data, &acc2); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(acc.AccountData, acc2.AccountData) {
		t.Errorf("matching error:\nneed: %+v\ngot:  %+v", acc, acc2)
	}
}
