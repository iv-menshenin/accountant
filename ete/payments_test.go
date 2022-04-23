package ete_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"testing"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func testPayments(t *testing.T, logData fmt.Stringer, actor httpActor) {
	var personID = uuid.NewUUID()
	var objectID = uuid.NewUUID()
	var accountID = uuid.NewUUID()
	var targetID = uuid.NewUUID()
	var testPaymentsData = []domain.Payment{
		{
			AccountID: accountID,
			PaymentData: domain.PaymentData{
				PersonID: &personID,
				ObjectID: &objectID,
				Period: domain.Period{
					Month: 9,
					Year:  2021,
				},
				Target: domain.TargetHead{
					TargetID: targetID,
					Type:     "test",
				},
				Payment:     12000,
				PaymentDate: nil,
				Receipt:     "",
			},
		},
		{
			AccountID: accountID,
			PaymentData: domain.PaymentData{
				PersonID: nil,
				ObjectID: nil,
				Period: domain.Period{
					Month: 10,
					Year:  2021,
				},
				Target: domain.TargetHead{
					TargetID: targetID,
					Type:     "test",
				},
				Payment:     3400,
				PaymentDate: nil,
				Receipt:     "",
			},
		},
		{
			AccountID: uuid.NewUUID(),
			PaymentData: domain.PaymentData{
				PersonID: nil,
				ObjectID: nil,
				Period: domain.Period{
					Month: 11,
					Year:  2021,
				},
				Target: domain.TargetHead{
					TargetID: targetID,
					Type:     "test",
				},
				Payment:     1230,
				PaymentDate: nil,
				Receipt:     "",
			},
		},
	}

	t.Run("fill", func(t *testing.T) {
		for i, payment := range testPaymentsData {
			pay, err := actor.createPayment(payment.AccountID, payment.PaymentData)
			if err != nil {
				t.Fatal(err)
			}
			testPaymentsData[i].PaymentID = pay.PaymentID
		}
	})

	t.Run("fill", func(t *testing.T) {
		for i, payment := range testPaymentsData {
			pay, err := actor.getPayment(payment.PaymentID)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(testPaymentsData[i], *pay) {
				t.Fatalf("matching error\nwant: %+v\n get: %+v", testPaymentsData[i], *pay)
			}
		}
	})

	t.Run("find", func(t *testing.T) {
		found, err := actor.findPayment(&targetID, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		sort.Slice(found, func(i, j int) bool {
			return found[i].Period.Year*12+found[i].Period.Month < found[j].Period.Year*12+found[j].Period.Month
		})
		if !reflect.DeepEqual(testPaymentsData, found) {
			t.Fatalf("matching error\nwant: %+v\n get: %+v", testPaymentsData, found)
		}

		found, err = actor.findPayment(nil, &accountID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		sort.Slice(found, func(i, j int) bool {
			return found[i].Period.Year*12+found[i].Period.Month < found[j].Period.Year*12+found[j].Period.Month
		})
		if !reflect.DeepEqual(testPaymentsData[:2], found) {
			t.Fatalf("matching error\nwant: %+v\n get: %+v", testPaymentsData[:2], found)
		}

		found, err = actor.findPayment(nil, &accountID, &personID, &objectID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(testPaymentsData[:1], found) {
			t.Fatalf("matching error\nwant: %+v\n get: %+v", testPaymentsData[:1], found)
		}
	})

	t.Run("delete", func(t *testing.T) {
		for _, payment := range testPaymentsData {
			err := actor.deletePayment(payment.PaymentID)
			if err != nil {
				t.Fatal(err)
			}
			_, err = actor.getPayment(payment.PaymentID)
			if err != ErrNotFound {
				t.Fatalf("expected ErrNotFound, got: %v", err)
			}
		}
	})
}

func (a *httpActor) createPayment(accountID uuid.UUID, data domain.PaymentData) (result *domain.Payment, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, a.accountsURL+"/"+accountID.String()+pathPayments, buf)
	if err != nil {
		return
	}
	var respData PaymentDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) getPayment(uuid uuid.UUID) (result *domain.Payment, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.paymentsURL+"/"+uuid.String(), nil)
	if err != nil {
		return
	}
	var respData PaymentDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) findPayment(targetID, accountID, personID, objectID *uuid.UUID) (result []domain.Payment, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.paymentsURL, nil)
	if err != nil {
		return
	}
	var values = url.Values{}
	if targetID != nil {
		values.Add("target_id", targetID.String())
	}
	if accountID != nil {
		values.Add("account_id", accountID.String())
	}
	if personID != nil {
		values.Add("person_id", personID.String())
	}
	if objectID != nil {
		values.Add("object_id", objectID.String())
	}
	req.URL.RawQuery = values.Encode()

	var respData PaymentsDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) deletePayment(uuid uuid.UUID) error {
	req, err := http.NewRequest(http.MethodDelete, a.paymentsURL+"/"+uuid.String(), nil)
	if err != nil {
		return err
	}
	var respData PaymentDataResponse
	if err = a.exec(req, &respData); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}
