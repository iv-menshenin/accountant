package ete_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/iv-menshenin/accountant/business"
	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/storage/mongodb"
	"github.com/iv-menshenin/accountant/transport"
)

func Test_ete(t *testing.T) {
	const (
		proto        = "http"
		port         = 8088
		host         = "localhost"
		pathAPI      = "/api"
		pathAccounts = "/accounts"
	)
	logData := bytes.NewBufferString("")
	var (
		listeningError = make(chan error)
		logger         = log.New(logData, "test logger::", 0)
		dbConfig       = config.New(
			"db",
			"-db.mongo-host=172.16.35.129",
			"-db.mongo-dbname=test",
			"-db.mongo-username=mongo",
			"-db.mongo-password=gfhjkm",
		)
		httpConfig = config.New(
			"http",
			fmt.Sprintf("-http.http-port=%d", port),
			fmt.Sprintf("-http.http-host=%s", host),
		)
		accountCreateURL = fmt.Sprintf("%s://%s:%d%s%s", proto, host, port, pathAPI, pathAccounts)
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

	var someDate = time.Date(2016, 12, 12, 0, 0, 0, 0, time.UTC)
	var buf = bytes.NewBufferString("")
	var account = model.AccountData{
		Account:       "000001",
		CadNumber:     "4535-34543-345343-34534",
		AgreementNum:  "NNN-0000001",
		AgreementDate: &someDate,
		PurchaseKind:  "договор",
		PurchaseDate:  time.Now(),
		Comment:       "test comment",
	}
	enc := json.NewEncoder(buf)
	if err = enc.Encode(account); err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest(http.MethodPost, accountCreateURL, buf)
	if err != nil {
		t.Log(logData.String())
		t.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log(logData.String())
		t.Fatal(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		t.Log(logData.String())
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if err = queryTransport.Shutdown(context.Background()); err != nil {
		t.Log(logData.String())
		t.Fatal(err)
	}

	select {
	case err = <-listeningError:
		t.Log(logData.String())
		t.Error(err)
	case <-time.After(time.Millisecond * 5):
		//ok
	}
}
