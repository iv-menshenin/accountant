package ete_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iv-menshenin/accountant/utils/uuid"
	"io"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/business"
	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/model/domain"
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
	pathTargets  = "/targets"
	pathBills    = "/bills"
	pathPayments = "/payments"

	mongoDbHost = "172.16.35.129"
	mongoDbName = "test"
	mongoDbUser = "mongo"
	mongoDbPass = "gfhjkm"
)

var ErrNotFound = errors.New("not found")

func Test_ete(t *testing.T) {
	logData := bytes.NewBufferString("")
	var actor = upService(t, logData)

	var (
		account1 *domain.Account
		account2 *domain.Account
	)

	t.Run("AccountTesting", func(t *testing.T) {
		account1 = testAccount(t, logData, actor, account1)
		account2 = testAccount(t, logData, actor, account2)
	})
	t.Run("PersonTesting", func(t *testing.T) {
		account1 = testPerson(t, logData, actor, account1)
		account2 = testPerson(t, logData, actor, account2)
	})
	t.Run("ObjectTesting", func(t *testing.T) {
		account1 = testObject(t, logData, actor, account1)
		account2 = testObject(t, logData, actor, account2)
	})

	t.Run("GetAccounts", func(t *testing.T) {
		_, err := actor.getAccounts("")
		if err != nil {
			t.Errorf("cannot get accounts: %s", err)
		}
		accs, err := actor.getAccounts("account=" + account1.Account)
		if err != nil {
			t.Errorf("cannot get accounts: %s", err)
		}
		if len(accs) != 1 || !accs[0].AccountID.Equal(account1.AccountID) {
			t.Errorf("unexpected result: %+v", accs)
		}
	})

	t.Run("Finalization", func(t *testing.T) {
		wipeAccount(t, logData, actor, account1)
		wipeAccount(t, logData, actor, account2)
	})

	t.Run("Targets", func(t *testing.T) {
		testTargets(t, logData, actor)
	})

	t.Run("Bills", func(t *testing.T) {
		testBills(t, logData, actor)
	})

	t.Run("Payments", func(t *testing.T) {
		testPayments(t, logData, actor)
	})

	if err := actor.release(); err != nil {
		t.Log(logData.String())
		t.Fatal(err)
	}
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
			fmt.Sprintf("-http.port=%d", port),
			fmt.Sprintf("-http.host=%s", host),
		)
		accountsURL = fmt.Sprintf("%s://%s:%d%s%s", proto, host, port, pathAPI, pathAccounts)
		targetsURL  = fmt.Sprintf("%s://%s:%d%s%s", proto, host, port, pathAPI, pathTargets)
		billsURL    = fmt.Sprintf("%s://%s:%d%s%s", proto, host, port, pathAPI, pathBills)
		paymentsURL = fmt.Sprintf("%s://%s:%d%s%s", proto, host, port, pathAPI, pathPayments)
	)
	mongoStorage, err := mongodb.NewStorage(dbConfig, logger)
	if err != nil {
		t.Fatal(err)
	}

	var (
		accountCollection  = mongoStorage.NewAccountsCollection(storage.MapMongodbErrors)
		personsCollection  = mongoStorage.NewPersonsCollection(accountCollection, storage.MapMongodbErrors)
		objectsCollection  = mongoStorage.NewObjectsCollection(accountCollection, storage.MapMongodbErrors)
		targetsCollection  = mongoStorage.NewTargetsCollection(storage.MapMongodbErrors)
		billsCollection    = mongoStorage.NewBillsCollection(storage.MapMongodbErrors)
		paymentsCollection = mongoStorage.NewPaymentsCollection(storage.MapMongodbErrors)

		appHnd = business.New(
			&testLogger{l: logger},
			accountCollection,
			personsCollection,
			objectsCollection,
			targetsCollection,
			billsCollection,
			paymentsCollection,
		)
		queryTransport = transport.NewHTTPServer(httpConfig, logger, appHnd, nil)
	)
	go queryTransport.ListenAndServe(listeningError)
	// todo health check
	time.Sleep(time.Millisecond * 100)

	return httpActor{
		accountsURL: accountsURL,
		targetsURL:  targetsURL,
		billsURL:    billsURL,
		paymentsURL: paymentsURL,
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
		targetsURL  string
		billsURL    string
		paymentsURL string
		release     func() error
	}
	ResponseMeta struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	}
	AccountsDataResponse struct {
		Meta ResponseMeta     `json:"meta"`
		Data []domain.Account `json:"data,omitempty"`
	}
	AccountDataResponse struct {
		Meta ResponseMeta    `json:"meta"`
		Data *domain.Account `json:"data,omitempty"`
	}
	PersonDataResponse struct {
		Meta ResponseMeta   `json:"meta"`
		Data *domain.Person `json:"data,omitempty"`
	}
	ObjectDataResponse struct {
		Meta ResponseMeta   `json:"meta"`
		Data *domain.Object `json:"data,omitempty"`
	}
	TargetDataResponse struct {
		Meta ResponseMeta   `json:"meta"`
		Data *domain.Target `json:"data,omitempty"`
	}
	TargetsDataResponse struct {
		Meta ResponseMeta    `json:"meta"`
		Data []domain.Target `json:"data,omitempty"`
	}
	BillDataResponse struct {
		Meta ResponseMeta `json:"meta"`
		Data *domain.Bill `json:"data,omitempty"`
	}
	BillsDataResponse struct {
		Meta ResponseMeta  `json:"meta"`
		Data []domain.Bill `json:"data,omitempty"`
	}
	PaymentDataResponse struct {
		Meta ResponseMeta    `json:"meta"`
		Data *domain.Payment `json:"data,omitempty"`
	}
	PaymentsDataResponse struct {
		Meta ResponseMeta     `json:"meta"`
		Data []domain.Payment `json:"data,omitempty"`
	}
)

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

func mustUUID(u string) (x uuid.UUID) {
	if err := x.FromString(u); err != nil {
		panic(err)
	}
	return
}
