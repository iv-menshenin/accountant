package mongodb

import (
	"context"
	"fmt"
	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/storage"
	"log"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

var testStorage *Storage
var once sync.Once

func initTestEnv() {
	rand.Seed(time.Now().UnixNano())
	var err error
	var logger = log.Default()
	testStorage, err = NewStorage(config.New("tst"), logger)
	if err != nil {
		panic(err)
	}
}

type mock struct {
	accountMock []model.Account
}

const (
	defaultMockSize = 128
	lightMockSize   = 32
)

func newMock(mockSize int) *mock {
	var accountMock []model.Account
	for nn := 0; nn < mockSize; nn++ {
		var account = makeAccount(nn)
		for i := 0; i < rand.Intn(3)+1; i++ {
			account.Persons = append(account.Persons, makePerson(nn))
		}
		for i := 0; i < rand.Intn(3)+1; i++ {
			account.Objects = append(account.Objects, makeObject(nn))
		}
		accountMock = append(accountMock, account)
	}
	return &mock{
		accountMock: accountMock,
	}
}

type accManipulator struct {
	*AccountCollection
	*PersonCollection
	*ObjectCollection
}

func (m *accManipulator) uploadAccount(ctx context.Context, acc model.Account) error {
	err := m.AccountCollection.Create(ctx, model.Account{
		AccountID:   acc.AccountID,
		Persons:     nil,
		Objects:     nil,
		AccountData: acc.AccountData,
	})
	if err != nil {
		return err
	}
	for _, person := range acc.Persons {
		err = m.PersonCollection.Create(ctx, acc.AccountID, person)
		if err != nil {
			return err
		}
	}
	for _, object := range acc.Objects {
		err = m.ObjectCollection.Create(ctx, acc.AccountID, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func Test_StorageAccounts(t *testing.T) {
	once.Do(initTestEnv)

	accounts := testStorage.NewAccountCollection(storage.MapMongodbErrors)
	persons := testStorage.NewPersonCollection(accounts, storage.MapMongodbErrors)
	objects := testStorage.NewObjectCollection(accounts, storage.MapMongodbErrors)

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
		go func(acc *model.Account) {
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
		go func(acc *model.Account) {
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

			accs, err := accounts.Find(ctx, model.FindAccountOption{Account: &acc.Account})
			if err != nil {
				errCh <- err
				return
			}
			if len(accs) != 1 || !reflect.DeepEqual(accs[0], *acc) {
				errCh <- fmt.Errorf("[ACCOUNT] want: %+v, got: %+v", *acc, accs)
				return
			}

			var foundPersons []model.Person
			foundPersons, err = persons.Find(ctx, model.FindPersonOption{
				AccountID: &acc.AccountID,
			})
			if !reflect.DeepEqual(acc.Persons, foundPersons) {
				errCh <- fmt.Errorf("[PERSONS] want: %+v, got: %+v", acc.Persons, foundPersons)
				return
			}

			var foundObjects []model.Object
			foundObjects, err = objects.Find(ctx, model.FindObjectOption{
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
			accs, err = accounts.Find(ctx, model.FindAccountOption{Account: &found.Account})
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

func Test_StorageObjects(t *testing.T) {
	once.Do(initTestEnv)

	accounts := testStorage.NewAccountCollection(storage.MapMongodbErrors)
	persons := testStorage.NewPersonCollection(accounts, storage.MapMongodbErrors)
	objects := testStorage.NewObjectCollection(accounts, storage.MapMongodbErrors)

	var manipulator = accManipulator{
		accounts,
		persons,
		objects,
	}
	var data = newMock(lightMockSize)
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
		go func(acc *model.Account) {
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
				found, err := objects.Find(ctx, model.FindObjectOption{
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
					itsOk = itsOk || reflect.DeepEqual(f, obj)
				}
				if !itsOk {
					errCh <- fmt.Errorf("cant find object by address: %v", obj)
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
			objs, err := objects.Find(ctx, model.FindObjectOption{AccountID: &acc.AccountID})
			if err != nil {
				errCh <- fmt.Errorf("cant find object: %w", err)
				return
			}
			if !reflect.DeepEqual(objs, acc.Objects) {
				errCh <- fmt.Errorf("error matching objects: %v, got: %v", acc.Objects, objs)
				return
			}

			err = objects.Delete(ctx, acc.AccountID, rndObj.ObjectID)
			if err != nil {
				errCh <- fmt.Errorf("cant delete object: %w", err)
				return
			}
			acc.Objects = append(acc.Objects[:rndObjNum], acc.Objects[rndObjNum+1:]...)
			objs, err = objects.Find(ctx, model.FindObjectOption{AccountID: &acc.AccountID})
			if err != nil {
				errCh <- fmt.Errorf("cant find object: %w", err)
				return
			}
			if !reflect.DeepEqual(objs, acc.Objects) {
				errCh <- fmt.Errorf("error matching objects: %v, got: %v", acc.Objects, objs)
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

func makeAccount(nn int) model.Account {
	var purchased = map[int]string{
		0: "договор",
		1: "дарение",
		2: "наследство",
		3: "украл",
	}
	agrDate := time.Date(rand.Intn(30)+1980, time.Month(rand.Intn(12)+1), rand.Intn(30)+1, 0, 0, 0, 0, time.UTC)
	purchDate := time.Date(rand.Intn(30)+1980, time.Month(rand.Intn(12)+1), rand.Intn(30)+1, 0, 0, 0, 0, time.UTC)
	return model.Account{
		AccountID: uuid.NewUUID(),
		AccountData: model.AccountData{
			Account:       fmt.Sprintf("#%d", rand.Int()),
			CadNumber:     fmt.Sprintf("%d:%d:%d:%d", rand.Intn(5)+81, rand.Intn(89)+10, rand.Intn(1999999)+1000000, nn),
			AgreementNum:  fmt.Sprintf("№ %d:%d", rand.Intn(5)+81, nn),
			AgreementDate: &agrDate,
			PurchaseKind:  purchased[nn%4],
			PurchaseDate:  purchDate,
			Comment:       fmt.Sprintf("test account rnd = %d", rand.Int()),
		},
	}
}

func makePerson(nn int) model.Person {
	var surnames = map[int]string{
		0: "Карасёв",
		1: "Дунаев",
		2: "Чебыкина",
		3: "Волянский",
		4: "Лукьяненко",
		5: "Асатрян",
	}
	var names = map[int]string{
		0: "Валентин",
		1: "Владимир",
		2: "Александр",
		3: "Василий",
		4: "Дмитрий",
		5: "Михаил",
	}
	var patnames = map[int]string{
		0: "Васильевич",
		1: "Павлович",
		2: "Валентинович",
		3: "Иванович",
		4: "Игоревич",
		5: "Борисович",
	}
	return model.Person{
		PersonID: uuid.NewUUID(),
		PersonData: model.PersonData{
			Name:     names[rand.Intn(len(names))],
			Surname:  surnames[rand.Intn(len(surnames))],
			PatName:  patnames[rand.Intn(len(patnames))],
			DOB:      nil,
			IsMember: nn%2 > 0,
			Phone:    fmt.Sprintf("8(922)-%d-%d-%d", rand.Intn(99)+100, rand.Intn(56)+10, rand.Intn(56)+10),
			EMail:    fmt.Sprintf("%d@mail.ru", rand.Int()),
		},
	}
}

func makeObject(nn int) model.Object {
	var streets = map[int]string{
		0: "Фруктовая",
		1: "Вишневая",
		2: "Сыромолотова",
		3: "Чайная",
		4: "Инженерная",
		5: "Осенняя",
		6: "Весенняя",
		7: "Зимняя",
		8: "Коричневая",
		9: "Сталинская",
	}
	return model.Object{
		ObjectID: uuid.NewUUID(),
		ObjectData: model.ObjectData{
			PostalCode: "889-009",
			City:       "Krasnodar",
			Village:    "Victoria",
			Street:     streets[rand.Intn(len(streets))],
			Number:     nn,
			Area:       float64(rand.Intn(20) + 390),
		},
	}
}
