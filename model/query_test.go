package model

import (
	"encoding/json"
	"github.com/iv-menshenin/accountant/model/uuid"
	"reflect"
	"testing"
	"time"
)

func Test_SaveAccountQuery(t *testing.T) {
	var acc = Account{
		AccountID: uuid.NewUUID(),
		AccountData: AccountData{
			Comment:      "foo, bar",
			AgreementNum: "#foo-001-002-00bar",
			AgreementDate: func() *time.Time {
				var tm = time.Now().Round(time.Hour)
				return &tm
			}(),
		},
	}
	var aq = PostAccountQuery{
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
