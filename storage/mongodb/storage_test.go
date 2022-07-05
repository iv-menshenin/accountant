package mongodb

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

var testStorage *Storage
var once sync.Once

func TestMain(m *testing.M) {
	once.Do(initTestEnv)
	defer func() {
		if e := testStorage.Close(); e != nil {
			panic(e)
		}
	}()
	os.Exit(m.Run())
}

func initTestEnv() {
	rand.Seed(time.Now().UnixNano())
	var err error
	var logger = log.Default()
	testStorage, err = NewTestStorage(config.New("tst"), logger)
	if err != nil {
		panic(err)
	}
}

type mock struct {
	accountMock []domain.Account
}

const (
	defaultMockSize = 128
)

func newMock(mockSize int) *mock {
	var accountMock []domain.Account
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
	*AccountsCollection
	*PersonsCollection
	*ObjectsCollection
}

func (m *accManipulator) uploadAccount(ctx context.Context, acc domain.Account) error {
	err := m.AccountsCollection.Create(ctx, domain.Account{
		AccountID:   acc.AccountID,
		Persons:     nil,
		Objects:     nil,
		AccountData: acc.AccountData,
	})
	if err != nil {
		return err
	}
	for _, person := range acc.Persons {
		err = m.PersonsCollection.Create(ctx, acc.AccountID, person)
		if err != nil {
			return err
		}
	}
	for _, object := range acc.Objects {
		err = m.ObjectsCollection.Create(ctx, acc.AccountID, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeAccount(nn int) domain.Account {
	var purchased = map[int]string{
		0: "договор",
		1: "дарение",
		2: "наследство",
		3: "украл",
	}
	agrDate := time.Date(rand.Intn(30)+1980, time.Month(rand.Intn(12)+1), rand.Intn(30)+1, 0, 0, 0, 0, time.UTC)
	purchDate := time.Date(rand.Intn(30)+1980, time.Month(rand.Intn(12)+1), rand.Intn(30)+1, 0, 0, 0, 0, time.UTC)
	return domain.Account{
		AccountID: uuid.NewUUID(),
		AccountData: domain.AccountData{
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

func makePerson(nn int) domain.Person {
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
	return domain.Person{
		PersonID: uuid.NewUUID(),
		PersonData: domain.PersonData{
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

func makeObject(nn int) domain.Object {
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
	return domain.Object{
		ObjectID: uuid.NewUUID(),
		ObjectData: domain.ObjectData{
			PostalCode: "889-009",
			City:       "Krasnodar",
			Village:    "Victoria",
			Street:     streets[rand.Intn(len(streets))],
			Number:     nn,
			Area:       float64(rand.Intn(20) + 390),
		},
	}
}
