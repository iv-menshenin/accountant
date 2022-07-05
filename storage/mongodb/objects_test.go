package mongodb

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func Test_Objects(t *testing.T) {
	once.Do(initTestEnv)

	accounts := testStorage.NewAccountsCollection(storage.MapMongodbErrors)
	persons := testStorage.NewPersonsCollection(accounts, storage.MapMongodbErrors)
	objects := testStorage.NewObjectsCollection(accounts, storage.MapMongodbErrors)

	var manipulator = accManipulator{
		accounts,
		persons,
		objects,
	}
	var data = newMock(defaultMockSize)
	var accountMock = data.accountMock
	var wg sync.WaitGroup
	var errCh = make(chan error)
	wg.Add(len(accountMock))

	var closed = make(chan struct{})
	go func() {
		defer close(closed)
		for err := range errCh {
			t.Error(err)
		}
	}()

	for i := range accountMock {
		go func(acc *domain.Account) {
			defer wg.Done()

			if len(acc.Objects) == 0 {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()

			err := manipulator.uploadAccount(ctx, *acc)
			if err != nil {
				errCh <- fmt.Errorf("cant upload account: %w", err)
				return
			}

			for _, obj := range acc.Objects {
				found, err := objects.Find(ctx, storage.FindObjectOption{
					AccountID: nil,
					Address:   &obj.Street,
					Number:    &obj.Number,
				})
				if err != nil {
					errCh <- fmt.Errorf("cant find objects by address: %w", err)
					return
				}
				var itsOk bool
				for _, f := range found {
					itsOk = itsOk || reflect.DeepEqual(f, domain.NestedObject{
						Object:    obj,
						AccountID: acc.AccountID,
					})
				}
				if !itsOk {
					errCh <- fmt.Errorf("cant find object by address: %v\n: %v", obj, found)
					return
				}
			}

			rndObjNum := rand.Intn(len(acc.Objects))
			rndObj := acc.Objects[rndObjNum]
			looked, err := objects.Lookup(ctx, acc.AccountID, rndObj.ObjectID)
			if err != nil {
				errCh <- fmt.Errorf("cant lookup objects by ID: %w", err)
				return
			}
			if looked == nil || !reflect.DeepEqual(rndObj, *looked) {
				errCh <- fmt.Errorf("cant lookup object by ID: want: %v, got: %v", rndObj, looked)
				return
			}

			rndObj.Street = "replaced"
			rndObj.Number = 9901
			rndObj.PostalCode = "000000"
			err = objects.Replace(ctx, acc.AccountID, rndObj.ObjectID, rndObj)
			if err != nil {
				errCh <- fmt.Errorf("cant replace objects by ID: %w", err)
				return
			}

			acc.Objects[rndObjNum] = rndObj
			objs, err := objects.Find(ctx, storage.FindObjectOption{AccountID: &acc.AccountID})
			if err != nil {
				errCh <- fmt.Errorf("cant find object: %w", err)
				return
			}
			var needObjs = objectsToNeed(acc.Objects, acc.AccountID)
			if !reflect.DeepEqual(objs, needObjs) {
				errCh <- fmt.Errorf("error matching objects\nneed: %v\n got: %v", needObjs, objs)
				return
			}

			err = objects.Delete(ctx, acc.AccountID, rndObj.ObjectID)
			if err != nil {
				errCh <- fmt.Errorf("cant delete object: %w", err)
				return
			}
			acc.Objects = append(acc.Objects[:rndObjNum], acc.Objects[rndObjNum+1:]...)
			objs, err = objects.Find(ctx, storage.FindObjectOption{AccountID: &acc.AccountID})
			if err != nil {
				errCh <- fmt.Errorf("cant find object: %w", err)
				return
			}
			needObjects := objectsToNeed(acc.Objects, acc.AccountID)
			if !reflect.DeepEqual(objs, needObjects) {
				errCh <- fmt.Errorf("error matching objects: %v, got: %v", needObjects, objs)
				return
			}

			err = objects.Delete(ctx, acc.AccountID, uuid.NewUUID())
			if err != storage.ErrNotFound {
				errCh <- fmt.Errorf("expected storage.ErrNotFound, got: %w", err)
				return
			}

		}(&accountMock[i])
	}

	wg.Wait()
	close(errCh)
	<-closed
}

func objectsToNeed(objects []domain.Object, accountID uuid.UUID) []domain.NestedObject {
	var needObjs = make([]domain.NestedObject, 0)
	for _, object := range objects {
		needObjs = append(needObjs, domain.NestedObject{
			Object:    object,
			AccountID: accountID,
		})
	}
	return needObjs
}
