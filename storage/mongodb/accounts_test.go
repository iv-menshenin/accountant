package mongodb

import (
	"context"
	"fmt"
	"github.com/iv-menshenin/accountant/utils/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
)

func Test_Accounts(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var wg sync.WaitGroup
	var errCh = make(chan error)

	wg.Add(len(accountMock))
	for i := range accountMock {
		go func(acc *domain.Account) {
			defer wg.Done()
			if err := manipulator.uploadAccount(ctx, *acc); err != nil {
				errCh <- err
			}
		}(&accountMock[i])
	}

	var closed = make(chan struct{})
	go func() {
		defer close(closed)
		for err := range errCh {
			t.Error(err)
		}
	}()
	wg.Wait()

	wg.Add(len(accountMock))
	for i := range accountMock {
		go func(acc *domain.Account) {
			defer wg.Done()

			found, err := accounts.Lookup(ctx, acc.AccountID)
			if err != nil {
				errCh <- err
				return
			}
			if !reflect.DeepEqual(acc, found) {
				errCh <- fmt.Errorf("[ACCOUNT] want: %+v, got: %+v", *acc, *found)
				return
			}

			accs, err := accounts.Find(ctx, storage.FindAccountOption{Account: &acc.Account})
			if err != nil {
				errCh <- err
				return
			}
			if len(accs) != 1 || !reflect.DeepEqual(accs[0], *acc) {
				errCh <- fmt.Errorf("[ACCOUNT] want: %+v, got: %+v", *acc, accs)
				return
			}

			var foundObjects []domain.NestedObject
			var needObjects = objectsToNeed(acc.Objects, acc.AccountID)
			foundObjects, err = objects.Find(ctx, storage.FindObjectOption{
				AccountID: &acc.AccountID,
			})
			if !reflect.DeepEqual(needObjects, foundObjects) {
				errCh <- fmt.Errorf("[OBJECTS] want: %+v, got: %+v", needObjects, foundObjects)
				return
			}

			if err = accounts.Delete(ctx, acc.AccountID); err != nil {
				errCh <- err
				return
			}
			accs, err = accounts.Find(ctx, storage.FindAccountOption{Account: &found.Account})
			if err != nil {
				errCh <- err
				return
			}
			if len(accs) != 0 {
				errCh <- fmt.Errorf("[ACCOUNT] must be deleted, but found: %+v", accs)
				return
			}

		}(&accountMock[i])
	}

	wg.Wait()
	close(errCh)
	<-closed
}

func Test_Find_Accounts(t *testing.T) {
	once.Do(initTestEnv)

	accounts := testStorage.NewAccountsCollection(storage.MapMongodbErrors)
	persons := testStorage.NewPersonsCollection(accounts, storage.MapMongodbErrors)
	objects := testStorage.NewObjectsCollection(accounts, storage.MapMongodbErrors)

	_, err := accounts.storage.DeleteMany(context.Background(), bson.M{"persons.name": "Ktulhu"})
	if err != nil {
		panic(err)
	}

	var manipulator = accManipulator{
		accounts,
		persons,
		objects,
	}
	var case1 = domain.Account{
		AccountID: uuid.NewUUID(),
		Persons: []domain.Person{
			{
				PersonID: uuid.NewUUID(),
				PersonData: domain.PersonData{
					Name:    "Ktulhu",
					Surname: "Fhtagn",
					PatName: "Vladimirovitch",
				},
			},
		},
		Objects: []domain.Object{
			{
				ObjectID: uuid.NewUUID(),
				ObjectData: domain.ObjectData{
					Street: "Ingeenernaya",
					Number: 54,
				},
			},
		},
	}
	var testCases = []domain.Account{
		case1,
	}
	for _, test := range testCases {
		if err := manipulator.AccountsCollection.Create(context.Background(), test); err != nil {
			t.Error(err)
		}
	}

	acc, err := manipulator.AccountsCollection.Find(context.Background(), storage.FindAccountOption{ByPerson: "tulhu"})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(acc, []domain.Account{case1}) {
		t.Errorf("matching error\nwant: %+v\n got: %+v", []domain.Account{case1}, acc)
	}
}
