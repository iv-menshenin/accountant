package ete_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/iv-menshenin/accountant/model/uuid"
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

func Test_ete(t *testing.T) {
	logData := bytes.NewBufferString("")
	var actor = upService(t, logData)

	var account *model.Account

	account = testAccount(t, logData, actor, account)
	account = testPerson(t, logData, actor, account)

	if err := actor.release(); err != nil {
		t.Log(logData.String())
		t.Fatal(err)
	}
}

func testAccount(t *testing.T, logData fmt.Stringer, actor httpActor, account *model.Account) *model.Account {
	t.Run("create user", func(t *testing.T) {
		var err error
		var someDate = time.Date(2016, 12, 12, 0, 0, 0, 0, time.UTC)
		var accountData = model.AccountData{
			Account:       fmt.Sprintf("%d", rand.Int()),
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
	return account
}

func testPerson(t *testing.T, logData fmt.Stringer, actor httpActor, account *model.Account) *model.Account {
	t.Run("create_person", func(t *testing.T) {

		var person = model.PersonData{
			Name:     "Test",
			Surname:  "Testing",
			PatName:  "Ivanovich",
			IsMember: true,
			Phone:    "(555)-555-55-55",
			EMail:    "test@accountant",
		}
		gotPerson, err := actor.createPerson(account.AccountID, person)
		if err != nil {
			t.Log(logData.String())
			t.Fatal(err)
		}

		account.Persons = append(account.Persons, *gotPerson)

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
		accountCollection = mongoStorage.NewAccountCollection(storage.MapError)
		personsCollection = mongoStorage.NewPersonCollection(accountCollection, storage.MapError)
		objectsCollection = mongoStorage.NewObjectCollection(accountCollection, storage.MapError)

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
	if resp.StatusCode > 399 {
		return fmt.Errorf("unexpected http status: %d", resp.StatusCode)
	}
	return nil
}
