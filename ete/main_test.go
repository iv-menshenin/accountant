package ete_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/business"
	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/storage/mongodb"
	"github.com/iv-menshenin/accountant/transport"
)

const (
	proto        = "http"
	port         = 8088
	host         = "localhost"
	pathAPI      = "/api"
	pathAccounts = "/accounts"

	mongoDbHost = "172.16.35.129"
	mongoDbName = "test"
	mongoDbUser = "mongo"
	mongoDbPass = "gfhjkm"
)

var ErrNotFound = errors.New("not found")

func Test_ete(t *testing.T) {
	logData := bytes.NewBufferString("")
	var actor = upService(t, logData)

	var account *model.Account

	t.Run("AccountTesting", func(t *testing.T) {
		account = testAccount(t, logData, actor, account)
	})
	t.Run("PersonTesting", func(t *testing.T) {
		account = testPerson(t, logData, actor, account)
	})
	t.Run("PersonObject", func(t *testing.T) {
		account = testObject(t, logData, actor, account)
	})
	t.Run("Finalization", func(t *testing.T) {
		wipeAccount(t, logData, actor, account)
	})

	if err := actor.release(); err != nil {
		t.Log(logData.String())
		t.Fatal(err)
	}
}

func testAccount(t *testing.T, logData fmt.Stringer, actor httpActor, account *model.Account) *model.Account {
	t.Run("create account", func(t *testing.T) {
		var err error
		var someDate = time.Date(2016, 12, 12, 0, 0, 0, 0, time.UTC)
		var accountData = model.AccountData{
			Account:       fmt.Sprintf("%d", rand.Intn(8999999)+1000000),
			CadNumber:     "4535-34543-345343-34534",
			AgreementNum:  "NNN-0000001",
			AgreementDate: &someDate,
			PurchaseKind:  "договор",
			PurchaseDate:  time.Date(2006, 6, 6, 0, 0, 0, 0, time.UTC),
			Comment:       "test comment",
		}
		account, err = actor.createAccount(accountData)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(accountData, account.AccountData) {
			t.Log(logData.String())
			t.Fatalf("error while account creating\nwant: %+v\n got: %+v", accountData, account.AccountData)
		}

		got, err := actor.getAccount(account.AccountID)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("error while account getting\nwant: %+v\n got: %+v", accountData, account.AccountData)
		} else if !reflect.DeepEqual(got, account) {
			t.Log(logData.String())
			t.Fatalf("error while account creating\nwant: %+v\n got: %+v", account, got)
		}
	})
	t.Run("update account", func(t *testing.T) {
		var err error
		var someDate = time.Date(2016, 12, 12, 0, 0, 0, 0, time.UTC)
		var accountData = model.AccountData{
			Account:       fmt.Sprintf("%d", rand.Intn(8999999)+1000000),
			CadNumber:     "4535-34543-345343-34534",
			AgreementNum:  fmt.Sprintf("№ %d", rand.Intn(89999)+10000),
			AgreementDate: &someDate,
			PurchaseKind:  "договор",
			PurchaseDate:  time.Date(2006, 6, 14, 0, 0, 0, 0, time.UTC),
			Comment:       "комментарий",
		}
		account, err = actor.updateAccount(account.AccountID, accountData)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(accountData, account.AccountData) {
			t.Log(logData.String())
			t.Fatalf("error while account updating\nwant: %+v\n got: %+v", accountData, account.AccountData)
		}

		got, err := actor.getAccount(account.AccountID)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("error while account getting\nwant: %+v\n got: %+v", accountData, account.AccountData)
		} else if !reflect.DeepEqual(got, account) {
			t.Log(logData.String())
			t.Fatalf("error while account creating\nwant: %+v\n got: %+v", account, got)
		}
	})
	return account
}

func wipeAccount(t *testing.T, logData fmt.Stringer, actor httpActor, account *model.Account) {
	t.Run("delete account", func(t *testing.T) {
		err := actor.deleteAccount(account.AccountID)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		}

		got, err := actor.getAccount(account.AccountID)
		if err == ErrNotFound {
			return
		}
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		}
		if got != nil {
			t.Log(logData.String())
			t.Fatal(fmt.Errorf("account has not been deleted: %+v", *got))
		}
	})
}

func testPerson(t *testing.T, logData fmt.Stringer, actor httpActor, account *model.Account) *model.Account {
	t.Run("create_person", func(t *testing.T) {

		var person = model.PersonData{
			Name:     "Test1",
			Surname:  "Testing",
			PatName:  "Ivanovich",
			IsMember: true,
			Phone:    "(555)-555-55-55",
			EMail:    "foo@accountant",
		}
		gotPerson, err := actor.createPerson(account.AccountID, person)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotPerson.PersonData, person) {
			t.Log(logData.String())
			t.Fatalf("error while person creating\nwant: %+v\n got: %+v", person, gotPerson.PersonData)
		}

		account.Persons = append(account.Persons, *gotPerson)
	})
	t.Run("add_new_person", func(t *testing.T) {
		var person = model.PersonData{
			Name:     "Test2",
			Surname:  "Testing",
			PatName:  "Michailovna",
			IsMember: true,
			Phone:    "(555)-777-77-77",
			EMail:    "bar@accountant",
		}
		gotPerson, err := actor.createPerson(account.AccountID, person)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotPerson.PersonData, person) {
			t.Log(logData.String())
			t.Fatalf("error while person creating\nwant: %+v\n got: %+v", person, gotPerson.PersonData)
		}

		account.Persons = append(account.Persons, *gotPerson)
	})
	t.Run("check_persons", func(t *testing.T) {
		got, err := actor.getAccount(account.AccountID)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("error while account getting\nwant: %+v\n", account.AccountData)
		} else if !reflect.DeepEqual(got, account) {
			t.Log(logData.String())
			t.Fatalf("error while account creating\nwant: %+v\n got: %+v", account, got)
		}
	})
	t.Run("delete_first_person", func(t *testing.T) {
		err := actor.deletePerson(account.AccountID, account.Persons[0].PersonID)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		}
		account.Persons = account.Persons[1:]
	})
	t.Run("update_person", func(t *testing.T) {
		var person = account.Persons[0]
		person.Name = "Измененный"
		person.Surname = "Васильев"
		person.PatName = "Веникович"
		person.IsMember = false
		account.Persons[0] = person

		gotPerson, err := actor.updatePerson(account.AccountID, account.Persons[0].PersonID, person.PersonData)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotPerson.PersonData, person.PersonData) {
			t.Log(logData.String())
			t.Fatalf("error while person updating\nwant: %+v\n got: %+v", person.PersonData, gotPerson.PersonData)
		}
	})
	t.Run("check_deleted_persons", func(t *testing.T) {
		got, err := actor.getAccount(account.AccountID)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("error while account getting\nwant: %+v\n", account.AccountData)
		} else if !reflect.DeepEqual(got, account) {
			t.Log(logData.String())
			t.Fatalf("error while account creating\nwant: %+v\n got: %+v", account, got)
		}
	})
	return account
}

func testObject(t *testing.T, logData fmt.Stringer, actor httpActor, account *model.Account) *model.Account {
	t.Run("create_object", func(t *testing.T) {

		var object = model.ObjectData{
			PostalCode: "663451",
			City:       "г. Краснокрыльск",
			Village:    "пос. Прозорливый",
			Street:     "ул. Ответственности",
			Number:     99,
			Area:       422.5,
		}
		gotObject, err := actor.createObject(account.AccountID, object)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotObject.ObjectData, object) {
			t.Log(logData.String())
			t.Fatalf("error while object creating\nwant: %+v\n got: %+v", object, gotObject.ObjectData)
		}

		account.Objects = append(account.Objects, *gotObject)
	})
	t.Run("add_new_object", func(t *testing.T) {
		var object = model.ObjectData{
			PostalCode: "",
			City:       "г. Краснокрыльск",
			Village:    "пос. Прозорливый",
			Street:     "ул. Ответственности",
			Number:     98,
			Area:       422.5,
		}
		gotObject, err := actor.createObject(account.AccountID, object)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotObject.ObjectData, object) {
			t.Log(logData.String())
			t.Fatalf("error while object creating\nwant: %+v\n got: %+v", object, gotObject.ObjectData)
		}

		account.Objects = append(account.Objects, *gotObject)
	})
	t.Run("add_new_object", func(t *testing.T) {
		var object = account.Objects[1]
		object.Number = 89
		object.PostalCode = "333423"
		account.Objects[1] = object
		gotObject, err := actor.updateObject(account.AccountID, object.ObjectID, object.ObjectData)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		} else if !reflect.DeepEqual(gotObject.ObjectData, object.ObjectData) {
			t.Log(logData.String())
			t.Fatalf("error while object updating\nwant: %+v\n got: %+v", object.ObjectData, gotObject.ObjectData)
		}
	})
	t.Run("check_objects", func(t *testing.T) {
		got, err := actor.getAccount(account.AccountID)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("error while account getting\nwant: %+v\n", account.AccountData)
		} else if !reflect.DeepEqual(got, account) {
			t.Log(logData.String())
			t.Fatalf("error while account creating\nwant: %+v\n got: %+v", account, got)
		}
	})
	t.Run("delete_object", func(t *testing.T) {
		err := actor.deleteObject(account.AccountID, account.Objects[0].ObjectID)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		}
		account.Objects = account.Objects[1:]
	})
	t.Run("delete_last_object", func(t *testing.T) {
		err := actor.deleteObject(account.AccountID, account.Objects[0].ObjectID)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		}
		account.Objects = account.Objects[1:]
	})
	t.Run("check_empty_objects", func(t *testing.T) {
		got, err := actor.getAccount(account.AccountID)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("error while account getting\nwant: %+v\n", account.AccountData)
		} else if !reflect.DeepEqual(got, account) {
			t.Log(logData.String())
			t.Fatalf("error while account creating\nwant: %+v\n got: %+v", account, got)
		}
	})
	return account
}

func upService(t *testing.T, logData io.Writer) httpActor {
	rand.Seed(time.Now().UnixNano())
	var (
		listeningError = make(chan error)
		logger         = log.New(logData, "test logger::", 0)
		dbConfig       = config.New(
			"db",
			"-db.mongo-host="+mongoDbHost,
			"-db.mongo-dbname="+mongoDbName,
			"-db.mongo-username="+mongoDbUser,
			"-db.mongo-password="+mongoDbPass,
		)
		httpConfig = config.New(
			"http",
			fmt.Sprintf("-http.http-port=%d", port),
			fmt.Sprintf("-http.http-host=%s", host),
		)
		accountsURL = fmt.Sprintf("%s://%s:%d%s%s", proto, host, port, pathAPI, pathAccounts)
	)
	mongoStorage, err := mongodb.NewStorage(dbConfig, logger)
	if err != nil {
		t.Fatal(err)
	}

	var (
		accountCollection = mongoStorage.NewAccountCollection(storage.MapMongodbErrors)
		personsCollection = mongoStorage.NewPersonCollection(accountCollection, storage.MapMongodbErrors)
		objectsCollection = mongoStorage.NewObjectCollection(accountCollection, storage.MapMongodbErrors)

		appHnd         = business.New(&testLogger{l: logger}, accountCollection, personsCollection, objectsCollection)
		queryTransport = transport.NewHTTPServer(httpConfig, logger, appHnd, nil)
	)
	go queryTransport.ListenAndServe(listeningError)
	// todo health check
	time.Sleep(time.Millisecond * 100)

	return httpActor{
		accountsURL: accountsURL,
		release: func() error {
			if err = queryTransport.Shutdown(context.Background()); err != nil {
				return err
			}
			select {
			case err, ok := <-listeningError:
				if ok {
					return err
				}
				return nil
			case <-time.After(time.Millisecond * 5):
				return nil
			}
		},
	}
}

type (
	httpActor struct {
		accountsURL string
		release     func() error
	}
	ResponseMeta struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	}
	AccountDataResponse struct {
		Meta ResponseMeta   `json:"meta"`
		Data *model.Account `json:"data,omitempty"`
	}
	PersonDataResponse struct {
		Meta ResponseMeta  `json:"meta"`
		Data *model.Person `json:"data,omitempty"`
	}
	ObjectDataResponse struct {
		Meta ResponseMeta  `json:"meta"`
		Data *model.Object `json:"data,omitempty"`
	}
)

func (a *httpActor) createAccount(data model.AccountData) (result *model.Account, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, a.accountsURL, buf)
	if err != nil {
		return
	}
	var respData AccountDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) updateAccount(accID uuid.UUID, data model.AccountData) (result *model.Account, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPut, a.accountsURL+"/"+accID.String(), buf)
	if err != nil {
		return
	}
	var respData AccountDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) getAccount(accID uuid.UUID) (result *model.Account, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.accountsURL+"/"+accID.String(), nil)
	if err != nil {
		return nil, err
	}
	var respData AccountDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) deleteAccount(accID uuid.UUID) error {
	var req *http.Request
	req, err := http.NewRequest(http.MethodDelete, a.accountsURL+"/"+accID.String(), nil)
	if err != nil {
		return err
	}
	var respData AccountDataResponse
	if err = a.exec(req, &respData); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return fmt.Errorf("unexpected status: %s %s", respData.Meta.Status, respData.Meta.Message)
}

func (a *httpActor) createPerson(accID uuid.UUID, data model.PersonData) (result *model.Person, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, a.accountsURL+"/"+accID.String()+"/persons", buf)
	if err != nil {
		return nil, err
	}
	var respData PersonDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) updatePerson(accID, personID uuid.UUID, data model.PersonData) (result *model.Person, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPut, a.accountsURL+"/"+accID.String()+"/persons/"+personID.String(), buf)
	if err != nil {
		return nil, err
	}
	var respData PersonDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) deletePerson(accID, personID uuid.UUID) error {
	var req *http.Request
	req, err := http.NewRequest(http.MethodDelete, a.accountsURL+"/"+accID.String()+"/persons/"+personID.String(), nil)
	if err != nil {
		return err
	}
	var respData PersonDataResponse
	if err = a.exec(req, &respData); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return fmt.Errorf("unexpected status: %s %s", respData.Meta.Status, respData.Meta.Message)
}

func (a *httpActor) createObject(accID uuid.UUID, data model.ObjectData) (result *model.Object, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, a.accountsURL+"/"+accID.String()+"/objects", buf)
	if err != nil {
		return nil, err
	}
	var respData ObjectDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) updateObject(accID, objID uuid.UUID, data model.ObjectData) (result *model.Object, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPut, a.accountsURL+"/"+accID.String()+"/objects/"+objID.String(), buf)
	if err != nil {
		return nil, err
	}
	var respData ObjectDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) deleteObject(accID, objID uuid.UUID) error {
	var req *http.Request
	req, err := http.NewRequest(http.MethodDelete, a.accountsURL+"/"+accID.String()+"/objects/"+objID.String(), nil)
	if err != nil {
		return err
	}
	var respData ObjectDataResponse
	if err = a.exec(req, &respData); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return fmt.Errorf("unexpected status: %s %s", respData.Meta.Status, respData.Meta.Message)
}

func (a *httpActor) exec(req *http.Request, i interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil && err == nil {
			err = e
		}
	}()
	if resp.StatusCode == 200 {
		dec := json.NewDecoder(resp.Body)
		return dec.Decode(i)
	}
	if resp.StatusCode == 204 {
		return io.EOF
	}
	if resp.StatusCode == 404 {
		return ErrNotFound
	}
	if resp.StatusCode > 399 {
		return fmt.Errorf("unexpected http status: %d", resp.StatusCode)
	}
	return nil
}
