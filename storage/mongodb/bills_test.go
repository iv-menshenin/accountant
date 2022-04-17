package mongodb

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func Test_Bills(t *testing.T) {
	once.Do(initTestEnv)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := testStorage.mongo.Bills().DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("cannot delete bills: %s", err)
	}

	var accountID = uuid.NewUUID()
	bills := testStorage.NewBillsCollection(storage.MapMongodbErrors)
	testPeriod := []model.Bill{
		{
			BillID:    uuid.NewUUID(),
			AccountID: accountID,
			PersonID:  nil,
			ObjectID:  nil,
			Period: model.Period{
				Month: 5,
				Year:  2021,
			},
			Target: model.TargetHead{
				TargetID: uuid.NewUUID(),
				Type:     "test",
			},
			Bill: 1244,
			Payments: []uuid.UUID{
				uuid.NewUUID(), uuid.NewUUID(), uuid.NewUUID(), uuid.NewUUID(),
			},
		},
		{
			BillID:    uuid.NewUUID(),
			AccountID: accountID,
			Period: model.Period{
				Month: 5,
				Year:  2021,
			},
			Target: model.TargetHead{
				TargetID: uuid.NewUUID(),
				Type:     "test",
			},
			Bill: 3700,
			Payments: []uuid.UUID{
				uuid.NewUUID(),
			},
		},
	}
	testBills := append([]model.Bill{
		{
			BillID:    uuid.NewUUID(),
			AccountID: accountID,
			Period: model.Period{
				Month: 6,
				Year:  2021,
			},
			Target: model.TargetHead{
				TargetID: uuid.NewUUID(),
				Type:     "test2",
			},
			Bill: 2300,
		},
	}, testPeriod...)

	sort.Slice(testBills, func(i, j int) bool {
		return testBills[i].BillID.String() < testBills[j].BillID.String()
	})

	var noiceID = uuid.NewUUID()
	if err = bills.Create(ctx, model.Bill{
		BillID:    noiceID,
		AccountID: uuid.NewUUID(),
		Period: model.Period{
			Month: 2,
			Year:  2022,
		},
		Target: model.TargetHead{
			TargetID: uuid.NewUUID(),
			Type:     "noice",
		},
		Bill: 4333,
	}); err != nil {
		t.Errorf("cannot create new bill: %s", err)
	}

	for _, bill := range testBills {
		if err = bills.Create(ctx, bill); err != nil {
			t.Errorf("cannot create new bill: %s", err)
		}
	}

	found, err := bills.FindByAccount(ctx, accountID)
	if err != nil {
		t.Errorf("cannot find bill: %s", err)
	}
	if !reflect.DeepEqual(found, testBills) {
		t.Errorf("matching error. want: %v, got: %v", testBills, found)
	}

	for _, bill := range testBills {
		found, err := bills.Lookup(ctx, bill.BillID)
		if err != nil {
			t.Errorf("cannot find bill: %s", err)
		}
		if !reflect.DeepEqual(found, &bill) {
			t.Errorf("matching error.\nwant: %v\n got: %v", bill, found)
		}
	}

	err = bills.Delete(ctx, noiceID)
	if err != nil {
		t.Errorf("cannot delete bill: %s", err)
	}
	look, err := bills.FindByPeriod(ctx, model.Period{
		Month: 2,
		Year:  2022,
	})
	if err != nil {
		t.Errorf("cannot check deleted bill: %s", err)
	}
	if look != nil {
		t.Errorf("expected nil, got: %v", look)
	}

	found, err = bills.FindByPeriod(ctx, model.Period{
		Month: 5,
		Year:  2021,
	})
	if err != nil {
		t.Errorf("cannot found bills by period: %s", err)
	}
	sort.Slice(testPeriod, func(i, j int) bool {
		return testPeriod[i].BillID.String() < testPeriod[j].BillID.String()
	})
	sort.Slice(found, func(i, j int) bool {
		return found[i].BillID.String() < found[j].BillID.String()
	})
	if !reflect.DeepEqual(found, testPeriod) {
		t.Errorf("matching error.\nwant: %v\n got: %v", testPeriod, found)
	}
}
