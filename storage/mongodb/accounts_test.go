package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/model/domain"
	storage2 "github.com/iv-menshenin/accountant/model/storage"
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

			accs, err := accounts.Find(ctx, storage2.FindAccountOption{Account: &acc.Account})
			if err != nil {
				errCh <- err
				return
			}
			if len(accs) != 1 || !reflect.DeepEqual(accs[0], *acc) {
				errCh <- fmt.Errorf("[ACCOUNT] want: %+v, got: %+v", *acc, accs)
				return
			}

			var foundPersons []domain.Person
			foundPersons, err = persons.Find(ctx, storage2.FindPersonOption{
				AccountID: &acc.AccountID,
			})
			if !reflect.DeepEqual(acc.Persons, foundPersons) {
				errCh <- fmt.Errorf("[PERSONS] want: %+v, got: %+v", acc.Persons, foundPersons)
				return
			}

			var foundObjects []domain.Object
			foundObjects, err = objects.Find(ctx, storage2.FindObjectOption{
				AccountID: &acc.AccountID,
			})
			if !reflect.DeepEqual(acc.Objects, foundObjects) {
				errCh <- fmt.Errorf("[OBJECTS] want: %+v, got: %+v", acc.Objects, foundObjects)
				return
			}

			if err = accounts.Delete(ctx, acc.AccountID); err != nil {
				errCh <- err
				return
			}
			accs, err = accounts.Find(ctx, storage2.FindAccountOption{Account: &found.Account})
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
