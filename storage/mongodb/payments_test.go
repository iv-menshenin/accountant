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

func Test_Payments(t *testing.T) {
	once.Do(initTestEnv)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := testStorage.mongo.Payments().DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("cannot delete payments: %s", err)
	}

	payments := testStorage.NewPaymentCollection(storage.MapMongodbErrors)
	personID1 := uuid.NewUUID()
	personID2 := uuid.NewUUID()
	accountID := uuid.NewUUID()

	testPayments := []model.Payment{
		{
			PaymentID: uuid.NewUUID(),
			AccountID: accountID,
			PersonID:  &personID1,
			ObjectID:  nil,
			Period: model.Period{
				Month: 06,
				Year:  2012,
			},
			Target: model.TargetHead{
				TargetID: uuid.NewUUID(),
				Type:     "Ordinary",
			},
			Payment:     3400,
			PaymentDate: nil,
			Receipt:     "#444",
		},
		{
			PaymentID: uuid.NewUUID(),
			AccountID: accountID,
			PersonID:  &personID2,
			ObjectID:  nil,
			Period: model.Period{
				Month: 05,
				Year:  2012,
			},
			Target: model.TargetHead{
				TargetID: uuid.NewUUID(),
				Type:     "Ordinary",
			},
			Payment:     3400,
			PaymentDate: nil,
			Receipt:     "#444",
		},
	}

	sort.Slice(testPayments, func(i, j int) bool {
		return testPayments[i].PaymentID.String() < testPayments[j].PaymentID.String()
	})

	if err = payments.Create(ctx, testPayments[0]); err != nil {
		t.Fatalf("cannot create payment: %s", err)
	}

	if err = payments.Create(ctx, testPayments[1]); err != nil {
		t.Fatalf("cannot create payment: %s", err)
	}

	found, err := payments.FindByAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("cannot find payments: %s", err)
	}
	if !reflect.DeepEqual(testPayments, found) {
		t.Fatalf("want: %v, got: %v", testPayments, found)
	}

	findIDs := []uuid.UUID{testPayments[0].PaymentID, testPayments[1].PaymentID}
	found, err = payments.FindByIDs(ctx, findIDs)
	if err != nil {
		t.Fatalf("cannot find payments: %s", err)
	}
	if !reflect.DeepEqual(testPayments, found) {
		t.Fatalf("matching error\nwant: %v\n got: %v", testPayments, found)
	}

	err = payments.Delete(ctx, testPayments[1].PaymentID)
	if err != nil {
		t.Fatalf("cannot delete payment: %s", err)
	}

	found, err = payments.FindByIDs(ctx, findIDs)
	if err != nil {
		t.Fatalf("cannot find payments: %s", err)
	}
	if !reflect.DeepEqual(testPayments[:1], found) {
		t.Fatalf("matching error\nwant: %v\n got: %v", testPayments[:1], found)
	}

}
