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

func Test_Persons(t *testing.T) {
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

			if len(acc.Persons) == 0 {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()

			err := manipulator.uploadAccount(ctx, *acc)
			if err != nil {
				errCh <- fmt.Errorf("cant upload account: %w", err)
				return
			}

			rndPerNum := rand.Intn(len(acc.Persons))
			rndPer := acc.Persons[rndPerNum]
			looked, err := persons.Lookup(ctx, acc.AccountID, rndPer.PersonID)
			if err != nil {
				errCh <- fmt.Errorf("cant lookup objects by ID: %w", err)
				return
			}
			if looked == nil || !reflect.DeepEqual(rndPer, *looked) {
				errCh <- fmt.Errorf("cant lookup person by ID: want: %v, got: %v", rndPer, looked)
				return
			}

			rndPer.Surname = "Menshenin"
			rndPer.Name = "Igor"
			rndPer.Phone = "000000"
			err = persons.Replace(ctx, acc.AccountID, rndPer.PersonID, rndPer)
			if err != nil {
				errCh <- fmt.Errorf("cant replace person by ID: %w", err)
				return
			}

			err = persons.Delete(ctx, acc.AccountID, rndPer.PersonID)
			if err != nil {
				errCh <- fmt.Errorf("cant delete person: %w", err)
				return
			}

			err = persons.Delete(ctx, acc.AccountID, uuid.NewUUID())
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

func personsToNeed(persons []domain.Person, accountID uuid.UUID) []domain.NestedPerson {
	var needPersons = make([]domain.NestedPerson, 0)
	for _, person := range persons {
		needPersons = append(needPersons, domain.NestedPerson{
			Person:    person,
			AccountID: accountID,
		})
	}
	return needPersons
}
