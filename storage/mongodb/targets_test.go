package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func Test_Targets(t *testing.T) {
	once.Do(initTestEnv)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := testStorage.mongo.Targets().DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("cannot delete targets: %s", err)
	}

	var closed = time.Now().Round(time.Second).UTC()
	targets := testStorage.NewTargetsCollection(storage.MapMongodbErrors)
	testTarget := domain.Target{
		TargetHead: domain.TargetHead{
			TargetID: uuid.NewUUID(),
			Type:     "TEST-T",
		},
		TargetData: domain.TargetData{
			Period: domain.Period{
				Month: 12,
				Year:  2021,
			},
			Closed:  &closed,
			Cost:    123000,
			Comment: "test",
		},
	}
	arrTarget := []domain.Target{testTarget}

	if err := targets.Create(ctx, testTarget); err != nil {
		t.Fatalf("cannot create target: %s", err)
	}

	for n := 1; n < 12; n++ {
		if err := targets.Create(ctx, domain.Target{
			TargetHead: domain.TargetHead{
				TargetID: uuid.NewUUID(),
				Type:     fmt.Sprintf("Noice-%d", n),
			},
			TargetData: domain.TargetData{
				Period: domain.Period{
					Month: n,
					Year:  2021,
				},
				Cost:    77764,
				Comment: "test",
			},
		}); err != nil {
			t.Fatalf("cannot create target: %s", err)
		}
	}

	for n := 0; n < 6; n++ {
		y := 2016 + n
		target := domain.Target{
			TargetHead: domain.TargetHead{
				TargetID: uuid.NewUUID(),
				Type:     fmt.Sprintf("Dec-%d", n),
			},
			TargetData: domain.TargetData{
				Period: domain.Period{
					Month: 12,
					Year:  y,
				},
				Cost:    30000,
				Comment: "test-payload",
			},
		}
		if err := targets.Create(ctx, target); err != nil {
			t.Fatalf("cannot create target: %s", err)
		}
		if y == 2021 {
			arrTarget = append(arrTarget, target)
		}
	}

	look, err := targets.Lookup(ctx, testTarget.TargetID)
	if err != nil {
		t.Fatalf("cannot lookup target: %s", err)
	}
	if look == nil {
		t.Fatalf("cannot lookup target: returned nil")
	}
	if !reflect.DeepEqual(testTarget, *look) {
		t.Fatalf("want: %v, got: %v", testTarget, *look)
	}

	found, err := targets.FindByPeriod(ctx, storage.FindTargetOption{
		ShowClosed: true,
		Month:      12,
		Year:       2021,
	})
	if err != nil {
		t.Fatalf("cannot find targets: %s", err)
	}
	if !reflect.DeepEqual(arrTarget, found) {
		t.Fatalf("cannot find targets: %v", found)
	}

	err = targets.Delete(ctx, arrTarget[0].TargetID)
	if err != nil {
		t.Fatalf("cannot delete targets: %s", err)
	}
	arrTarget = arrTarget[1:]

	found, err = targets.FindByPeriod(ctx, storage.FindTargetOption{
		Month: 12,
		Year:  2021,
	})
	if err != nil {
		t.Fatalf("cannot find targets: %s", err)
	}
	if !reflect.DeepEqual(arrTarget, found) {
		t.Fatalf("cannot find targets: %v", found)
	}
}
