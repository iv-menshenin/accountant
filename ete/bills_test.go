package ete_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func testBills(t *testing.T, logData fmt.Stringer, actor httpActor) {
	var bills = []domain.Bill{
		{
			AccountID: mustUUID("bf4b910d-65bb-4eab-bafa-97f1eb37e376"),
			BillData: domain.BillData{
				Formed: time.Now().Round(time.Second).UTC(),
				Period: domain.Period{Month: 1, Year: 2021},
				Target: domain.TargetHead{
					TargetID: mustUUID("4bbbb2e9-acc5-436e-a7be-b9f665f80622"),
					Type:     "test1",
				},
				Bill: 123000,
			},
		},
		{
			AccountID: mustUUID("c898d046-be35-11ec-9d64-0242ac120002"),
			BillData: domain.BillData{
				Formed: time.Now().Round(time.Second).UTC(),
				Period: domain.Period{Month: 1, Year: 2021},
				Target: domain.TargetHead{
					TargetID: mustUUID("4bbbb2e9-acc5-436e-a7be-b9f665f80622"),
					Type:     "test1",
				},
				Bill: 89000,
			},
		},
		{
			AccountID: mustUUID("c898d046-be35-11ec-9d64-0242ac120002"),
			BillData: domain.BillData{
				Formed: time.Now().Round(time.Second).UTC(),
				Period: domain.Period{Month: 2, Year: 2021},
				Target: domain.TargetHead{
					TargetID: mustUUID("93bcccf9-34d3-4eed-ba78-9c64676f7e57"),
					Type:     "test2",
				},
				Bill: 5670,
			},
		},
		{
			AccountID: mustUUID("c898d046-be35-11ec-9d64-0242ac120002"),
			BillData: domain.BillData{
				Formed: time.Now().Round(time.Second).UTC(),
				Period: domain.Period{Month: 3, Year: 2021},
				Target: domain.TargetHead{
					TargetID: mustUUID("55f084a2-6903-4678-aaef-280ed0a9712a"),
					Type:     "test3",
				},
				Bill: 12000,
			},
		},
	}
	t.Run("create", func(t *testing.T) {
		for i, bill := range bills {
			b, err := actor.createBill(bill.AccountID, bill.BillData)
			if err != nil {
				t.Fatalf("can not create bill: %s", err)
			}
			bills[i].BillID = b.BillID
		}
	})
	t.Run("get", func(t *testing.T) {
		for _, bill := range bills {
			b, err := actor.getBill(bill.BillID)
			if err != nil {
				t.Fatalf("can not lookup bill: %s", err)
			}
			if !reflect.DeepEqual(*b, bill) {
				t.Errorf("matching error:\nwant: %v,\n got: %v", bill, *b)
			}
		}
	})
	t.Run("find", func(t *testing.T) {
		b, err := actor.findBillsByTargetID(mustUUID("4bbbb2e9-acc5-436e-a7be-b9f665f80622"))
		if err != nil {
			t.Fatalf("can not lookup bill: %s", err)
		}
		if !reflect.DeepEqual(b, bills[:2]) {
			t.Errorf("matching error:\nwant: %v,\n got: %v", bills[:2], b)
		}

		b, err = actor.findBillsByPeriod(domain.Period{Month: 1, Year: 2021})
		if err != nil {
			t.Fatalf("can not lookup bill: %s", err)
		}
		if !reflect.DeepEqual(b, bills[:2]) {
			t.Errorf("matching error:\nwant: %v,\n got: %v", bills[:2], b)
		}

		b, err = actor.findBillsByAccountID(mustUUID("c898d046-be35-11ec-9d64-0242ac120002"))
		if err != nil {
			t.Fatalf("can not lookup bill: %s", err)
		}
		if !reflect.DeepEqual(b, bills[1:]) {
			t.Errorf("matching error:\nwant: %v,\n got: %v", bills[1:], b)
		}
	})

	t.Run("delete", func(t *testing.T) {
		for _, bill := range bills {
			if err := actor.deleteBill(bill.BillID); err != nil {
				t.Fatalf("can not delete bill: %s", err)
			}
		}
		b, err := actor.findBillsByAccountID(mustUUID("c898d046-be35-11ec-9d64-0242ac120002"))
		if err != nil && err != ErrNotFound {
			t.Fatalf("can not lookup bill: %s", err)
		}
		if len(b) > 0 {
			t.Errorf("expected empty, got: %v", b)
		}
	})
}

func (a *httpActor) createBill(accountID uuid.UUID, data domain.BillData) (result *domain.Bill, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, a.accountsURL+"/"+accountID.String()+pathBills, buf)
	if err != nil {
		return
	}
	var respData BillDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) getBill(uuid uuid.UUID) (result *domain.Bill, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.billsURL+"/"+uuid.String(), nil)
	if err != nil {
		return
	}
	var respData BillDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) findBillsByPeriod(period domain.Period) (result []domain.Bill, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.billsURL, nil)
	if err != nil {
		return
	}
	req.URL.RawQuery = url.Values{
		"period": []string{fmt.Sprintf("%d.%d", period.Month, period.Year)},
	}.Encode()
	var respData BillsDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) findBillsByAccountID(accID uuid.UUID) (result []domain.Bill, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.billsURL, nil)
	if err != nil {
		return
	}
	req.URL.RawQuery = url.Values{
		"account_id": []string{accID.String()},
	}.Encode()
	var respData BillsDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) findBillsByTargetID(targetID uuid.UUID) (result []domain.Bill, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.billsURL, nil)
	if err != nil {
		return
	}
	req.URL.RawQuery = url.Values{
		"target_id": []string{targetID.String()},
	}.Encode()
	var respData BillsDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) deleteBill(uuid uuid.UUID) error {
	req, err := http.NewRequest(http.MethodDelete, a.billsURL+"/"+uuid.String(), nil)
	if err != nil {
		return err
	}
	var respData BillDataResponse
	if err = a.exec(req, &respData); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}
